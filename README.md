# Consolidated Job Portal
This web application gathers all the job openings from multiple job platform such as Glassdoor, Linkedin, etc.
This web application also allows users to contribute their job openings.

# Technology Used
Backend: Golang to get job openings and insert into Mongodb. 

Database:Mongodb

Front-end: HTML and CSS.

Infrastracture: This web application is deployed to the kubernetes cluster, which is constructed on AWS-EKS service.
              Mongodb data is volumed by creating PersistentVolume and statefulset.
              
DNS(Domain Network Service): DNS is done by AWS's Route53. The loadbalancer's IP is mapped to a record in Route53. 

Unit test: Test scripts are written by Go.

CI/CD Pipeline: Workflow is created by Github Actions.
