## howto build
# - become root (e.g. sudo su)
# - VHS_NO_SANDBOX="true" vhs < k3i.tape && chmod 666 k3i.gif

# Where should we write the GIF?
Output k3i.gif

Set FontSize 15
Set Width 1200
Set Height 750

# Start
Type "k3i"
Sleep 800ms
Enter
Sleep 5s

# Clear
Hide
Type "clear"
Enter
Show

# List
Type "k3i list"
Sleep 500ms
Enter
Sleep 3s

# Create help
Type "k3i create --help"
Sleep 500ms
Enter
Sleep 3s

# Create
Type "k3i create demo"
Sleep 500ms
Enter
Sleep 25s

# List
Type "k3i list"
Sleep 500ms
Enter
Sleep 3s

# shell
Type "k3i shell demo demo-w-0"
Sleep 500ms
Enter
Sleep 7s
Type "ifconfig eth0"
Sleep 500ms
Enter
Sleep 3s
Type "exit"
Sleep 500ms
Enter
Sleep 2s

# kubeconfig
Type "k3i kubeconfig demo --write"
Sleep 500ms
Enter
Sleep 3s

# Export kubeconfig
Type "export KUBECONFIG=./kubeconfig.demo.yaml"
Sleep 500ms
Enter
Sleep 3s

# Cluster info
Type "kubectl cluster-info"
Sleep 500ms
Enter
Sleep 3s

# Cluster nodes
Type "kubectl get nodes"
Sleep 500ms
Enter
Sleep 3s

# Wait
Type "Lets wait a few seconds for the cluster to become ready"
Sleep 5s
Backspace 55

# Cluster nodes
Type "kubectl get nodes -o wide"
Sleep 500ms
Enter
Sleep 5s

# # delete help
# Type "k3i delete --help"
# Sleep 500ms
# Enter
# Sleep 7s
# 
# # delete
# Type "k3i delete"
# Sleep 500ms
# Enter
# Sleep 7s
# 
# # List
# Type "k3i list"
# Sleep 500ms
# Enter
# Sleep 3s
