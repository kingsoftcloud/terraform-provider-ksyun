#!/bin/zsh

# run debug with goland and dlv
GOEXEC=$(which go)
DLVEXEC=/Applications/GoLand.app/Contents/plugins/go-plugin/lib/dlv/mac/dlv
WORKSTATION=/Users/chenbin/Documents/ksc_project/terraform-provider-ksyun
TEST_DIR=$WORKSTATION/__test

if [ ! -d "${TEST_DIR}" ]; then
  echo "test direction is not exist, will create it."
  mkdir "${TEST_DIR}"
fi

if [ $# -lt 1 ]; then
  echo "the debug function name must be delivered"
  exit 1
fi

if [ "$1" == "test" ]; then
  echo "the current environment is test env"
  export TF_ACC=true
  export KSYUN_ACCESS_KEY=AKLTd0hEEqwASUGa6TwtRqYQ
  export KSYUN_SECRET_KEY=OIMOlZeifHAjPnKKuGXKygacg4dSeh0CaNAY87TENnUSjM02GG3Vn9GoJU7FjP4g
fi

if [ "$1" == "online" ]; then
  echo "the current environment is online env"
  export TF_ACC=true
  export KSYUN_ACCESS_KEY=AKLTRFeltNhdQVSXafs74PUxqQ
  export KSYUN_SECRET_KEY=OCzZ7suA4ret9RABAsW2LYxJPLfS0cM38pwiWh6ebcueJrcDb8cgSPQSSL7HfoCPyQ
fi

DEBUGGER_FUNC_NAME=$2
DEBUGGER_FILE=$TEST_DIR/___$DEBUGGER_FUNC_NAME.test
DEBUGGER_OUTPUT=$DEBUGGER_FILE
DEBUGGER_LOG_FILE=$TEST_DIR/$DEBUGGER_FUNC_NAME.log

echo "compiling..."
$GOEXEC test -c -gcflags=-N -o $DEBUGGER_OUTPUT -gcflags all=-N /Users/chenbin/Documents/ksc_project/terraform-provider-ksyun/ksyun

echo "dlv running..."
exec $DLVEXEC --listen=:5566 --headless=true --api-version=2 --check-go-version=false --only-same-user=false \
 exec $DEBUGGER_OUTPUT -- -test.v -test.paniconexit0 -test.run $DEBUGGER_FUNC_NAME