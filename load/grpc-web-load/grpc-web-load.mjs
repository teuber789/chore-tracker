#!/usr/bin/env node

// This is an abomination, but it works for this project.
// XMLHttpRequest is needed because grpc-web is intended for use in browsers only.
// https://stackoverflow.com/a/77047149
// https://stackoverflow.com/a/59214301
import { createRequire } from "module";
const require = createRequire(import.meta.url)
global.XMLHttpRequest = require('xhr2');

import pb from './chore_tracker_pb.js';
import grpc from './chore_tracker_grpc_web_pb.js';
import { randomBytes } from "node:crypto";

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
async function doLoadTest() {
    try {
        // IRL, address and port would be env consts
        const client = new grpc.ChoreTrackerClient('http://localhost:8080');
        var count = 0

        // Run forever. Process will be SIGKILLed by controller.
        while (true) {
            // TODO remove this
            await new Promise(resolve => setTimeout(resolve, 500));

            // Add a family
            const family = await addFamily(client);
            console.log('Added family', family.getId(), family.getName());
            
            // Add children
            const child1 = await addChild(client, family.getId(), 1);
            console.log('Added child', child1.getId(), 'to family', child1.getFamilyId());
            const child2 = await addChild(client, family.getId(), 2);
            console.log('Added child', child2.getId(), 'to family', child2.getFamilyId());
            const child3 = await addChild(client, family.getId(), 3);
            console.log('Added child', child3.getId(), 'to family', child3.getFamilyId());

            // Create chore
            const chore1 = await createChore(client, family, 1);
            console.log('Added chore', chore1.getId(), 'to family', chore1.getFamilyId());
            const chore2 = await createChore(client, family, 2);
            console.log('Added chore', chore2.getId(), 'to family', chore2.getFamilyId());
            const chore3 = await createChore(client, family, 3);
            console.log('Added chore', chore3.getId(), 'to family', chore3.getFamilyId());
            
            // Get all chores
            await getAllChores(client, family);

            // Each child completes each chore
            const children = [child1, child2, child3];
            const chores = [chore1, chore2, chore3];
            for (const child of children) {
                for (const chore of chores) {
                    await markChoreCompleted(client, child, chore);
                }
            }

            // Each child retrieves all the chores they have done
            for (const child of children) {
                await getCompletedChores(client, child);
            }

            count++
            console.log('COUNT IS ' + count)
        }
        
    } catch (err) {
        console.error(err);
        process.exit = 1;
        return;
    }
}

doLoadTest();
