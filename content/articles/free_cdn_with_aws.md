## How to set up a free CDN for your self-hosted website?

### Why would you want to do this?

As you may, or may not, know, this website is self-hosted on my home kubernetes cluster (as of December 2022) which is underpowered to say the least. While the home cluster has been an amazing learning opportunity, and I have intentionally avoided service providers to not give myself the easy way out, there are still something that's better handled by them to keep the user experience enjoyable and snappy.

For example, pictures!

![drawing of the Abarat Archipelago](https://d3gwcypbkjc953.cloudfront.net/images/floating-nature.jpg)

I love this picture, I don't even know what it really is, other than what I gather from a [Google image search](https://onepiecefanon.fandom.com/wiki/The_Abarat_Archipelago?file=The_Abarat_Archipelago.jpg) (It's from a manga series called 'One Piece' and that it's publicly available for fair use under the CC license). If I wanted to share this photo, I could store it in the database or somewhere behind a web server like [NGINX](https://www.nginx.com/), and have my k8 cluster serve it up everytime I want to look at the photo. Sure, but there's a couple problems with that.

1. My home cluster is not beefy.
    * It's 5 [Pine64 Rock64's](https://www.pine64.org/devices/single-board-computers/rock64/) I got for $20 a pop a long time ago, each with only 2Gb of precious RAM.
2. My home cluster is just that, in my home.
    * Someone looking at my site (unlikely) on the West Coast will likely sprout a couple gray hairs before this picture loads.
3. Life is too short to not try new things.

### The Solution, Amazon Web Services.
So instead of rocking the roll-your-own method, for once we're going to try AWS by creating a free CDN for our website.

By the end of this, we'll have a fast, distributed method to deliver static content to support our site.

#### *BEWARE! AWS operates on a pay-as-you-go model, you could be left with a big bill if you don't monitor your usage and set up a budget through the AWS console*

### How to
#### 1. Register with AWS
head over to https://aws.amazon.com/console/ and sign up for an account.

* Heads up, you'll need a credit card to associate with account during the registration process.

#### 2. Activate MFA (optional)
To protect your account, I recommend setting up multifactor authentication (MFA) via the security credentials page.

### 3. Create an S3 bucket
Amazon described their Simple Storage Service (S3) as such

> an object storage service that offers industry-leading scalability, data availability, security, and performance.

To over simplify, just think of S3 as the place we'll actually store our content. that's its only job.

* From the AWS console homepage, bring up the `Services` menu from the top left. Scroll down to `Storage`, and select `S3`.
* Click `Create Bucket` and enter a Bucket Name. The name must be unique amongst all of AWS, and you should choose a [DNS-compliant name](https://duckduckgo.com/?t=ffab&q=DNS+compliant+name&ia=web) since this name will be part of the URI to access the bucket.
* Choose a region that's close to you for now.
* Scroll down to the card titled `Block Public Access settings for this bucket`
    * Uncheck `block all public access`
    * Check the acknowledgement
    * ponder the risks

### 4. Setup AWS CloudFront
At this point, we've created a way to store or static content, but we don't have a good way to access it.

Enter [AWS CloudFront](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/Introduction.html)

* From the AWS console homepage again,
    * Select the `Services` dropdown
    * Scroll down to `Networking & Content Delivery`
    * Click `CloudFront`
* On the next page, take a moment to read the AWS free tier limits and select `Create a CloudFront Distribution`
    * From the `Origin domain` select your recently made s3 bucket
    * Scroll down, click the radio boxes next to...
        * `redirect HTTP to HTTPS`
    * You can leave the remaining defaults for now.
    * Scroll to the bottom to Create Distribution
    * Copy the `Distribution domain name` to your new CloudFront and save it for later.

### 5. Create an IAM User/Group
At this point, our CDN is set up, but we need to create a way to access our content via [AWS' Identity and Access Management (IAM) service](https://docs.aws.amazon.com/IAM/latest/UserGuide/introduction.html).

To quote AWS...

> AWS Identity and Access Management (IAM) is a web service that helps you securely control access to AWS resources. You use IAM to control who is authenticated (signed in) and authorized (has permissions) to use resources.

* From the AWS console homepage again,
    * Select the `Services` dropdown
    * Scroll down to `Security, Identity, & Compliance` in the left menu.
    * Click `IAM`
* After page load, select `user` from the left menu, and click `Add User`
* Create a new IAM user
    * Add a username
    * Click the check next to `Access key - programatic access`
    * click next
* Create a new group and add our new IAM user
    * If you're familiar with *nix groups, these act just like that.
    * give the group a descriptive name
    * use the filter bar to add two security policies to your new group `AmazonS3FullAccess` & `CloudFrontFullAccess`
    * Add the new user to the group and create the user
* After you successfully create the user, a page will show your new user's Access ID and Key.
    * **Copy these and store them in a safe place. They will not be visible after leaving this page**.

### 6. Install the AWS CLI
Before we add content to our freshly configured S3 bucket, we are going to install the AWS Command Line Interface.

Check out the AWS CLI documentation for the most [up-to-date instructions on installing the AWS CLI on your OS](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html).

Check the cli installation
```shell
$ aws --version
```

Next, we need to log in our CLI client using the Access key and Private key we copied earlier. To do this, we’ll run aws configure:
```shell
$ aws configureAWS Access Key ID [None]: YOUR_ACCESS_KEYAWS Secret Access Key [None]: YOUR_SECRET_KEY Default region name [None]: e.g: us-east-1 Default output format [None]: LEAVE_BLANK
```
If nothing goes wrong, you should be logged in as your new IAM user with access to your swanky new S3 bucket.

### 7. Sync your Images
It's time to upload our content, lets sync a local directory to our S3 bucket.

We can use the `aws s3 sync` to keep a local directory synced with our s3 bucket
```shell
$ aws s3 sync /path/to/dir/with/images/ s3://[YOUR_S3_BUCKET_NAME]/images --delete --exclude ".*"
```
Replace the path to your directory and the name of you s3 bucket followed by the path you'd like to be to access the content (e.g., `/images`).

### 8. Congrats
That's it! You should now be able to access<sup>1</sup> your content at the cloudfront URL that we copied earlier. For us, it was https://d1738nvaauqo1n.cloudfront.net, so if we wanted to access our images/floating_nature.jpg from our fresh CDN, we would navigate to https://d1738nvaauqo1n.cloudfront.net/images/floating_nature.jpg

### 9. Troubleshooting AccessDenied errors
After initially creating my CDN, I received an AccessDenied Error in XML format. To get around this, navigate to your new S3 bucket, go to the Permissions Tab, and scroll down to the Bucket Policy section and create a new policy in JSON format.
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "PublicReadGetObject",
            "Effect": "Allow",
            "Principal": "*",
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::dpgraham-sss-test/*"
        }
    ]
}
```
replacing dpgraham-sss-test with your bucket's name.

 <sup>1</sup>Note, it may take some time before you can actually view your content. I've heard as short as 25 minutes, but it may take longer.
