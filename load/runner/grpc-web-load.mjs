#!/usr/bin/env node

// This is an abomination, but it works for this project.
// XMLHttpRequest is needed because grpc-web is intended for use in browsers only.
// https://stackoverflow.com/a/77047149
// https://stackoverflow.com/a/59214301
import { createRequire } from 'module';
const require = createRequire(import.meta.url)
global.XMLHttpRequest = require('xhr2');

import pb from './chore_tracker_pb.js';
import grpc from './chore_tracker_grpc_web_pb.js';
import { randomBytes } from 'node:crypto';

// Helper function that wraps client methods so their callbacks can be used as promises.
// Necessary because grpc-web only has limited support for TypeScript + promises.
function asPromise(f, req, metadata) {
    return new Promise((resolve, reject) => {
        f(req, metadata, (err, resp) => {
            if (err) {
                reject(err);
            } else {
                resolve(resp);
            }
        });
    });
}

function addFamily(client) {
    const req = new pb.AddFamilyRequest();
    const familyName = randomBytes(8).toString('hex');
    req.setName(familyName);
    return asPromise(client.addFamily.bind(client), req, {});
}

function addChild(client, familyId, num) {
    const age = (Math.floor(Math.random() * 17)) + 1;
    const req = new pb.AddChildRequest();
    req.setFamilyId(familyId);
    req.setName(`Family ${familyId} Child ${num}`);
    req.setAge(age);
    return asPromise(client.addChild.bind(client), req, {});
}

function createChore(client, family, num) {
    const name = `Family ${family.getId()} Chore ${num}`;
    const price = Math.random().toFixed(2);
    const req = new pb.CreateChoreRequest();
    req.setFamilyId(family.getId());
    req.setName(name);
    req.setDescription(name);
    req.setPrice(price);
    return asPromise(client.createChore.bind(client), req, {});
}

function getAllChores(client, family) {
    const pageable = new pb.Pageable();
    pageable.setPageToken("0");
    pageable.setPageSize(100);
    const req = new pb.GetChoresRequest();
    req.setPageable(pageable);
    req.setFamilyId(family.getId());
    req.setChildId(1);  // This isn't actually used by the service
    return asPromise(client.getChores.bind(client), req, {});
}

function markChoreCompleted(client, child, chore) {
    const req = new pb.MarkChoreCompletedRequest();
    req.setFamilyId(child.getFamilyId());
    req.setChildId(child.getId());
    req.setChoreId(chore.getId());
    return asPromise(client.markChoreCompleted.bind(client), req, {});
}

function getCompletedChores(client, child) {
    const pageable = new pb.Pageable();
    pageable.setPageToken("0");
    pageable.setPageSize(100);
    const req = new pb.GetChoresRequest();
    req.setPageable(pageable);
    req.setFamilyId(child.getFamilyId());
    req.setChildId(child.getId());
    return asPromise(client.getCompletedChores.bind(client), req, {});
}

// Runs the load test sequence repeatedly until SIGKILL occurs.
// TODO Simplify and clean up
async function doLoadTest() {
    const host = process.argv[2];

    try {
        // IRL, address and port would be env vars
        const client = new grpc.ChoreTrackerClient(`http://${host}:8080`);

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
                family = await addFamily(client);
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
                child1 = await addChild(client, family.getId(), 1);
                addChildSucceeded++;
                child2 = await addChild(client, family.getId(), 2);
                addChildSucceeded++;
                child3 = await addChild(client, family.getId(), 3);
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
                chore1 = await createChore(client, family, 1);
                createChoreSucceeded++;
                chore2 = await createChore(client, family, 2);
                createChoreSucceeded++;
                chore3 = await createChore(client, family, 3);
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
                await getAllChores(client, family);
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
                        await markChoreCompleted(client, child, chore);
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
                    await getCompletedChores(client, child);
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
