const {GetChoresRequest, GetChoresResponse, Pageable} = require('./chore_tracker_pb.js');
const {ChoreTrackerClient} = require('./chore_tracker_grpc_web_pb.js');

// IRL, address and port would be env vars
var client = new ChoreTrackerClient('http://localhost:8080');

var pageable = new Pageable();
pageable.setPageToken("0");
pageable.setPageSize(10);
var request = new GetChoresRequest();
request.setPageable(pageable);
request.setChildId(1);

client.getChores(request, {}, (err, response) => {
  if (err) {
    console.log('THERE WAS AN ERROR ' + err);
  } else {
    console.log(response.getChoresList()); 
  }
});
