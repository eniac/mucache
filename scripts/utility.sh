# Read log of a pod given the service name
# it removes the _ from the service name (e.g. my_service -> myservice)
function klog() {
  local service
  if [ -z "$1" ]; then
    echo "Please provide a service name"
    return 1
  fi
  service=${1//_/}
  echo "Getting logs for service: $service"
  kubectl get pods | grep "$service" | awk '{print $1}' | xargs -I {} kubectl logs {}
}

# Get the IP address of a service
function kip() {
  local service
  if [ -z "$1" ]; then
    echo "Please provide a service name"
    return 1
  fi
  service=${1//_/}
  kubectl get svc "$service" | tail -n 1 | awk '{print $3}'
}

function klogp() {
  local service
  if [ -z "$1" ]; then
    echo "Please provide a service name"
    return 1
  fi
  service=${1//_/}
  echo "Getting logs for service: $service"
  kubectl get pods | grep "$service" | awk '{print $1}' | xargs -I {} kubectl logs {} -p
}

# Get all ips of a service
function kips() {
  local service
  if [ -z "$1" ]; then
    echo "Please provide a service name"
    return 1
  fi
  service=${1//_/}
  kubectl get svc | grep LoadBalancer | grep "$service" | awk '{print $3}' | tr '\n' ' '
}
