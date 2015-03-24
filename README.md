drmaa2
======
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/dgruber/drmaa2)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://raw.githubusercontent.com/dgruber/drmaa2/master/COPYING)

A Go (#golang) API for job submission, job workflow management, and HPC cluster monitoring based on the open OGF DRMAA2 standard (http://www.drmaa.org / http://www.ogf.org).

This Go API is a wrapper around the DRMAA2 C APIs. It is not yet finished. Improvements, bug reports, and pull requests are welcome. DRMAA2 Go API is not yet standardized but a draft (based on this implementation) is available here: http://redmine.ogf.org/projects/drmaav2-go-binding/repository


## Installation

1. Download / Install Univa Grid Engine (see http://www.univa.com)
   Univa Grid Engine is currently the only cluster scheduler which supports
   the DRMAA2 standard. A quick guide to get it installed you can find in 
   my blog article: 
   http://www.gridengine.eu/index.php/grid-engine-internals/209-installing-a-univa-grid-engine-demo-cluster-with-one-command-2014-07-13

   Note that Univa Grid Engine is commercial.

2. Be sure that you have Grid Engine in your path (i.e. $SGE_ROOT is set)

   You can test this by calling qhost / qstat / .....
   Usually this is done by calling
    source /path/where/grid/engine/is/installed/default/common/settings.sh

   Note that "default" is the CELL directory which can be different to "default".

3. Clone the project with *git* tools in your go source path (github.com/dgruber/drmaa2).

4. Call 
    build.sh

## Usage

Please check out a detailed example which exploits the monitoring session and creates a
webserver hosting the cluster status information.

https://github.com/dgruber/drmaa2_web_app

For a simple multi-clustering (access multiple compute clusters with the same tools) application
based on this DRMAA2 API, please check out:

https://github.com/dgruber/ubercluster

It implements a simple webservice API for DRMAA2.

## Documentation

https://godoc.org/github.com/dgruber/drmaa2
