<!DOCTYPE html>
<html>

<body>

<h1>ChainVote Backend</h1>

<p>The ChainVote backend is a distributed blockchain system designed for the election process. It uses various technologies and methodologies for efficient and secure operations.</p>

<h2>Technologies Used</h2>
<ul>
  <li><strong>Message Broker:</strong> MQTT for high latency and low bandwidth distributed systems.</li>
  <li><strong>Consensus Algorithm:</strong> Custom Raft Algorithm implementation.</li>
  <li><strong>Data Management:</strong> Redis for state management.</li>
  <li><strong>Deployment:</strong> Docker for containerized deployment.</li>
</ul>

<h2>Architecture Overview</h2>
<p>The system utilizes MQTT as a message broker, chosen for its properties and compatibility with low network coverage. The Raft algorithm is employed for consensus, ensuring leader election and synchronization among nodes. Key points:</p>
<ul>
  <li>Nodes elect a leader, syncing logs and managing state changes.</li>
  <li>Transaction pool holds 5 transactions before processing into blocks.</li>
  <li>Each node maintains an independent SQLITE database for transactions.</li>
  <li>Blockchain files undergo verification before being added to the blockchain.</li>
  <li>HTTP servers with an administrator panel for node management and system control.</li>
</ul>

<h2>Administrator Panel</h2>
<p>The Vuejs-powered administrator panel offers functionalities for:</p>
<ul>
  <li>Creating counties, constituencies, wards, polling stations, candidates, and voters.</li>
  <li>Tracking registered desktop clients, distributed nodes, and transactions.</li>
  <li>Viewing blockchain details, including tallying results and per-county breakdowns.</li>
</ul>

<h2>User Access</h2>
<p>Access and operations within the administrator panel are role-based, restricted to authorized users based on their roles.The default user is admin@superuser.com and the password is 123456</p>

<h2>Additional Details</h2>
<ul>
  <li>Blockchain explorer for public viewing at http://localhost/block-explorer</li>
  <li>Each node has its own block explorer,voters access the block explorer on the current leader node</li>
  <li>Voter access to ballots in the blockchain using their transaction IDs received via SMS.</li>
  <li>Docker used for simulating different nodes, running on Alpine Linux containers.</li>
</ul>

<h2>Setup and Deployment</h2>
<p>The system is deployed using Docker containers, each simulating different nodes. The containers run on Alpine Linux.</p>

<h2>Frontend Interaction</h2>
<p>The frontend client at <a href="https://github.com/mark-judah/chainvote_dapp_client" target="_blank">ChainVote Dapp Client</a> makes requests to this backend for system interactions.</p>


</body>
</html>
