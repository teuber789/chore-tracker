#!/usr/bin/env node

import { randomBytes } from 'node:crypto';

// IRL, these would be env vars
const baseUrl = "http://127.0.0.1:8081";

async function addFamily() {
    const familyName = randomBytes(8).toString('hex');
    const data = {
        name: familyName,
    };

    const res = await fetch(`${baseUrl}/families`, {
        method: "POST",
        body: JSON.stringify(data),
    });
    return res.json();
}

async function addChild(family, num) {
    const age = (Math.floor(Math.random() * 17)) + 1;
    const data = {
        family_id: family.id,
        name: `Family ${family.id} Child ${num}`,
        age: age,
    };
    
    const res = await fetch(`${baseUrl}/children`, {
        method: "POST",
        body: JSON.stringify(data),
    });
    return res.json();
}

async function createChore(family, num) {
    const name = `Family ${family.id} Chore ${num}`;
    const price = Number(Math.random().toFixed(2));
    const data = {
        family_id: family.id,
        name: name,
        description: name,
        price: price,
    };

    const res = await fetch(`${baseUrl}/chores`, {
        method: "POST",
        body: JSON.stringify(data),
    });
    return res.json();
}

async function getAllChores(family) {
    const params = new URLSearchParams({
        pageToken: "0",
        pageSize: 100,
        familyId: family.id,
        childId: 1,  // This isn't actually used by the service
    });
    const url = `${baseUrl}/chores?${params.toString()}`;

    const res = await fetch(url);
    return res.json();
}

function markChoreCompleted(child, chore) {
    const data = {
        family_id: child.family_id,
        child_id: child.id,
        chore_id: chore.id,
    };

    return fetch(`${baseUrl}/completions`, {
        method: "POST",
        body: JSON.stringify(data),
    });
}

async function getCompletedChores(child) {
    const params = new URLSearchParams({
        pageToken: "0",
        pageSize: 100,
        familyId: child.family_id,
        childId: child.id,
    });
    const url = `${baseUrl}/chores?${params.toString()}`;

    const res = await fetch(url);
    return res.json();
}

// Runs the load test sequence repeatedly until SIGKILL occurs.
// TODO Simplify and clean up
async function doLoadTest() {
    try {
        // Metrics
        let itersFailed = 0;
        let itersSucceeded = 0;
        let addFamilyFailed = 0;
        let addFamilySucceeded = 0;
        let addChildFailed = 0;
        let addChildSucceeded = 0;
        let createChoreFailed = 0;
        let createChoreSucceeded = 0;
        let getAllChoresFailed = 0;
        let getAllChoresSucceeded = 0;
        let markChoreCompletedFailed = 0;
        let markChoreCompletedSucceeded = 0;
        let getCompletedChoresFailed = 0;
        let getCompletedChoresSucceeded = 0;

        // Run forever. Process can be stopped by interrupt or SIGKILL.
        while (true) {
            // Print output periodically
            if ((itersSucceeded + itersFailed) % 10 == 0) {
                console.log(
                    `Metrics: ` +
                    `itersSucceeded ${itersSucceeded} itersFailed ${itersFailed} ` +
                    `addFamilySucceeded ${addFamilySucceeded} addFamilyFailed ${addFamilyFailed} ` +
                    `addChildSucceeded ${addChildSucceeded} addChildFailed ${addChildFailed} ` +
                    `createChoreSucceeded ${createChoreSucceeded} createChoreFailed ${createChoreFailed} ` +
                    `getAllChoresSucceeded ${getAllChoresSucceeded} getAllChoresFailed ${getAllChoresFailed} ` + 
                    `markChoreCompletedSucceeded ${markChoreCompletedSucceeded} markChoreCompletedFailed ${markChoreCompletedFailed} ` +
                    `getCompletedChoresSucceeded ${getCompletedChoresSucceeded} getCompletedChoresFailed ${getCompletedChoresFailed}`
                );
            }

            // Add a family
            var family;
            try {
                family = await addFamily();
                addFamilySucceeded++;
            } catch (err) {
                console.log('Request to add family failed');
                console.error(err);
                itersFailed++;
                addFamilyFailed++;
                continue;
            }
            
            // Add children
            var child1;
            var child2;
            var child3;
            try {
                child1 = await addChild(family, 1);
                addChildSucceeded++;
                child2 = await addChild(family, 2);
                addChildSucceeded++;
                child3 = await addChild(family, 3);
                addChildSucceeded++;
            } catch(err) {
                console.log('Request to add child failed');
                console.error(err);
                itersFailed++;
                addChildFailed++;
                continue;
            }

            // Create chores
            var chore1;
            var chore2;
            var chore3;
            try {
                chore1 = await createChore(family, 1);
                createChoreSucceeded++;
                chore2 = await createChore(family, 2);
                createChoreSucceeded++;
                chore3 = await createChore(family, 3);
                createChoreSucceeded++;
            } catch(err) {
                console.log('Request to create chore failed');
                console.error(err);
                itersFailed++;
                createChoreFailed++;
                continue;
            }
            
            // Get all chores
            try {
                await getAllChores(family);
                getAllChoresSucceeded++;
            } catch(err) {
                console.log('Request to get all chores failed');
                console.error(err);
                itersFailed++;
                getAllChoresFailed++;
                continue;
            }

            // Each child completes each chore
            const children = [child1, child2, child3];
            const chores = [chore1, chore2, chore3];
            for (const child of children) {
                for (const chore of chores) {
                    try {
                        await markChoreCompleted(child, chore);
                        markChoreCompletedSucceeded++;
                    } catch(err) {
                        console.log('Request to mark chore completed failed');
                        console.error(err);
                        itersFailed++;
                        markChoreCompletedFailed++;
                        continue;
                    }
                }
            }

            // Each child retrieves all the chores they have done
            for (const child of children) {
                try {
                    await getCompletedChores(child);
                    getCompletedChoresSucceeded++;
                } catch(err) {
                    console.log('Request to get completed chores failed');
                    console.error(err);
                    itersFailed++;
                    getCompletedChoresFailed++;
                    continue;
                }
            }

            itersSucceeded++;
        }
        
    } catch (err) {
        console.error(err);
        process.exit = 1;
        return;
    }
}

doLoadTest();

