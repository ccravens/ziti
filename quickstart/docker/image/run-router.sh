#!/bin/bash

# give the controller time to ramp up before running if running in docker-compose
sleep 5

. "${ZITI_SCRIPTS}/ziti-cli-functions.sh"

if [[ "${ZITI_CTRL_EDGE_ADVERTISED_ADDRESS-}" == "" ]]; then export ZITI_CTRL_EDGE_ADVERTISED_ADDRESS="ziti-edge-controller"; fi
if [[ "${ZITI_EDGE_ROUTER_NAME-}" == "" ]]; then
  ZITI_EDGE_ROUTER_DESIRED_NAME="${ZITI_NETWORK-}-edge-router"
else
  ZITI_EDGE_ROUTER_DESIRED_NAME="${ZITI_EDGE_ROUTER_NAME}"
fi
if [[ "${ZITI_EDGE_ROUTER_PORT-}" == "" ]]; then export ZITI_EDGE_ROUTER_PORT="3022"; fi
if [[ "${ZITI_EDGE_ROUTER_HOSTNAME}" == "" ]]; then export ZITI_EDGE_ROUTER_HOSTNAME="${ZITI_EDGE_ROUTER_NAME}${ZITI_DOMAIN_SUFFIX}"; fi
if [[ "${ZITI_EDGE_ROUTER_ROLES}" == "" ]]; then export ZITI_EDGE_ROUTER_ROLES="${ZITI_EDGE_ROUTER_NAME}"; fi

. ${ZITI_HOME}/ziti.env

# This is a unique situation due to the env file being shared among all routers so we need to explicitly set the router
export ZITI_EDGE_ROUTER_NAME="${ZITI_EDGE_ROUTER_DESIRED_NAME}"

_wait_for_controller

zitiLogin

if [[ "$1" == "edge" ]]; then
  echo "CREATING EDGE ROUTER CONFIG"
  createEdgeRouterConfig "${ZITI_EDGE_ROUTER_NAME}"
fi
if [[ "$1" == "wss" ]]; then
  echo "CREATING EDGE ROUTER WSS CONFIG"
  createEdgeRouterWssConfig "${ZITI_EDGE_ROUTER_NAME}"
fi
if [[ "$1" == "fabric" ]]; then
  echo "CREATING FABRIC ROUTER CONFIG"
  createFabricRouterConfig "${ZITI_EDGE_ROUTER_NAME}"
fi
if [[ "$1" == "private" ]]; then
  echo "CREATING PRIVATE ROUTER CONFIG"
  createPrivateRouterConfig "${ZITI_EDGE_ROUTER_NAME}"
fi

echo "----------  Creating edge-router ${ZITI_EDGE_ROUTER_HOSTNAME}...."
found=$(ziti edge list edge-routers 'name = "'"${ZITI_EDGE_ROUTER_HOSTNAME}"'"' | grep -c "${ZITI_EDGE_ROUTER_HOSTNAME}")
if [[ found -gt 0 ]]; then
  echo "----------  Found existing edge-router ${ZITI_EDGE_ROUTER_HOSTNAME}...."
else
  "${ZITI_BIN_DIR}/ziti" edge create edge-router "${ZITI_EDGE_ROUTER_HOSTNAME}" -o "${ZITI_HOME}/${ZITI_EDGE_ROUTER_HOSTNAME}.jwt" -t -a "${ZITI_EDGE_ROUTER_ROLES}"
  sleep 1
  echo "---------- Enrolling edge-router ${ZITI_EDGE_ROUTER_HOSTNAME}...."
  "${ZITI_BIN_DIR}/ziti-router" enroll "${ZITI_HOME}/${ZITI_EDGE_ROUTER_HOSTNAME}.yaml" --jwt "${ZITI_HOME}/${ZITI_EDGE_ROUTER_HOSTNAME}.jwt"
  echo ""
fi

"${ZITI_BIN_DIR}/ziti-router" run "${ZITI_HOME}/${ZITI_EDGE_ROUTER_NAME}.yaml" > "${ZITI_HOME}/ziti-${ZITI_EDGE_ROUTER_HOSTNAME}.log"

