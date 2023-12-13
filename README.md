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
<p>Access and operations within the administrator panel are role-based, restricted to authorized users based on their roles.The default user is admin@superuser.com and the password is 123456. To access a nodes admin panel, go to http://localhost:8080</p>

<h2>Additional Details</h2>
<ul>
  <li>Blockchain explorer for public viewing at http://localhost/block-explorer</li>
  <li>Each node has its own block explorer,voters access the block explorer on the current leader node</li>
  <li>Voter access to ballots in the blockchain using their transaction IDs received via SMS.</li>
  <li>Docker used for simulating different nodes, running on Alpine Linux containers.</li>
</ul>

<h2>Setting Up the Backend</h2>

<h3>Clone the Project and Install Dependencies</h3>
<p>To start, clone the project and install Golang dependencies:</p>
<pre><code>
git clone https://github.com/your-username/ChainVoteBackend.git
cd ChainVoteBackend
# Install Golang dependencies
go mod download
</code></pre>

<h3>Install EMQX and Redis</h3>
<p>Install EMQX (MQTT broker) and Redis:</p>
<pre><code>
# Example for Ubuntu/Debian
sudo apt-get install emqx redis-server
</code></pre>

<h3>Configure MQTT Controller</h3>
<p>Update the MQTT broker URL in the <code>mqttController.go</code> file:</p>
<pre><code>
// Update MQTT broker URL
opts := mqtt.NewClientOptions().AddBroker("tcp://YOUR_MQTT_BROKER_URL:PORT")
</code></pre>

<h3>Linux Deployment</h3>
<p>Instructions for Linux:</p>
<ol>
  <li>Install EMQX and Redis as mentioned above.</li>
  <li>Clone the project, install Golang dependencies.</li>
  <li>Configure the MQTT controller with your broker URL.</li>
</ol>

<h2>Docker Setup</h2>

<h3>Building and Running Docker Image</h3>
<ol>
  <li><strong>Open a terminal or command prompt</strong>.</li>
  <li><strong>Navigate to the directory containing the Dockerfile</strong>:</li>
  <pre><code>
  cd /path/to/your/directory
  </code></pre>
  <li><strong>Run the Docker build command</strong>:</li>
  <pre><code>
  docker build -t your_image_name:tag .
  </code></pre>
</ol>

<p>Replace <code>your_image_name:tag</code> with your desired image name and tag.</p>

<p>Build and run the Docker image using Docker Compose:</p>
<pre><code>
# Build and run Docker image
docker-compose up
</code></pre>

<h3>Testing on a Single PC</h3>
<p>If testing on a single PC:</p>
<p>Edit the <code>docker-compose.yml</code> file and change the service name each time you run <code>docker-compose up</code>.</p>

<h2>Frontend Interaction</h2>
<p>The frontend client at <a href="https://github.com/mark-judah/chainvote_dapp_client" target="_blank">ChainVote Dapp Client</a> makes requests to this backend for system interactions.</p>


</body>
</html>
