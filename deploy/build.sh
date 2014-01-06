#!/bin/bash
if [ $# -lt 2 ]
then
  echo "build packageName appDir [configFileName]"
  exit 1
fi
if [ ! -d $2 ]
then
  echo "appDir '$2' not exists"
  exit 1
fi

WORKDIR=${PWD}
APPDIR=$2

export GOPATH=$GOPATH:$WORKDIR
echo "GOPATH = $GOPATH"
echo "building $1"
go install $1
if [ $? -ne 0 ]
then
  exit $?
fi
APPNAME=${1##*/}
echo "cp ${WORKDIR}/bin/$APPNAME >> $APPDIR/"
cp ${WORKDIR}/bin/$APPNAME $APPDIR/
if [ $? -ne 0 ]
then
  exit $?
fi
chmod 775 $APPDIR/$APPNAME

if [ $# -lt 3 ]
then
  exit 0
fi
CFGFILE=$3

if [ ! -f $WORKDIR/bin/config/$CFGFILE ]
then
  echo "config file '$WORKDIR/bin/config/$CFGFILE' not exists"
  exit 1
fi

if [ ! -d $APPDIR/config ]
then
  mkdir $APPDIR/config
fi
if [ ! -f $WORKDIR/config/$CFGFILE ]
then
  echo "cp $WORKDIR/bin/config/$CFGFILE >> $APPDIR/config/$CFGFILE"
  cp $WORKDIR/bin/config/$CFGFILE $APPDIR/config/$CFGFILE
else
  echo "SKIP cp $WORKDIR/bin/config/$CFGFILE >> $APPDIR/config/$CFGFILE"
fi