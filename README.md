# PMCaresDonationTracking
**Blockchain Based application to track PM Care fund**

This repository contains the code of a sample blockchain based application to raising and managing PM Care fund.

### Flow Diagram

![](images/architecture.png)

### Steps

Please follow the steps to run this application.

* Clone the code
* Setup IBP extension on VS Code
* Get connection profile and copy to `local` directory provided in github repository.
* To start the server, perform the following steps.

  ```
  cd backend
  npm install
  PORT=30001 DEPLOY_TYPE=local npm start
  ```
* Run the webapp.
  ```
  cd webapp
  npm install
  npm start
  ```
* Access the application

You can access the webapp at `http://localhost:3000`.
