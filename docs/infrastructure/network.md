# Notes on Network Infrastructure

This project uses an internal terraform module called `network` to create a network and related resources.

## Resources

### VPC

At the top of our network hierarchy is the [VPC](https://cloud.google.com/vpc/docs).

from the
GCP [google_compute_subnetwork terraform docs](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_subnetwork.html):

> Each GCP project contains one or more VPC networks. Each VPC network is a global entity spanning all GCP regions. This
> global VPC network allows VM instances and other resources to communicate with each other via internal, private IP
> addresses.

Not all services are created within the VPC, for example Cloud Run, Cloud Functions, and App Engine (weirdly enough). In
those instances you need to create a [connector](https://cloud.google.com/vpc/docs/configure-serverless-vpc-access).

GCP projects come with a default VPC, I do not use the default VPC and instead create a (custom, currently not in auto
mode) VPC.

Dev and Prod reside in separate VPCs. Ideally I would like to separate projects, and even billing accounts, but that is
overkill for this project. For now, we keep the two projects on different virtual networks.

ToDo: study up on differences in Global, Regional and Zonal resources.

### Subnets

>
> Each VPC network is subdivided into subnets, and each subnet is contained within a single region. You can have more
> than one subnet in a region for a given VPC network. Each subnet has a contiguous private RFC1918 IP space.

Subnets are a subdivision of a VPC. They allow us to isolate resources and for route traffic within a
VPC. For example, we can create a subnet (within a specified [CIDR](https://aws.amazon.com/what-is/cidr/) range) for our
databases

### Private Service Access

Here's GCP explanation
of how [Private Service Access and Private IP work together](https://cloud.google.com/sql/docs/postgres/private-ip):

> Configuring a Cloud SQL instance to use private IP requires private services access. Private services access lets you
> create private connections between your VPC network and the underlying Google service producer's VPC network. Google
> entities that offer services, such as Cloud SQL, are called service producers. Each Google service creates a subnet in
> which to provision resources. The subnet's IP address range is typically a /24 CIDR block that is chosen by the
> service
> and comes from the allocated IP address range.
>
> Private connections make services reachable without going through the internet or using external IP addresses. For
> this
> reason, private IP provides lower network latency than public IP.

An added benefit accessing our database through an internal (private) IP address is we essentially are set up
end to end for encryption in transit (assuming we use an SSL certificate with the global load balancer) from database to
the user's browser.