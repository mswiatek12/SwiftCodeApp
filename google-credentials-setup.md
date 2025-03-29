# How to Obtain Google API Credentials

Follow these steps to generate the necessary Google API credentials for interacting with Google Sheets.

## Step 1: Go to Google Cloud Console
Visit the [Google Cloud Console](https://console.cloud.google.com/).

## Step 2: Create a Project
1. Click on the **Select a Project** dropdown at the top of the page.
2. Click **New Project**.
3. Give the project a name, for example, "SwiftCodeApp".
4. Click **Create** to create your project.

## Step 3: Enable APIs
1. In the left sidebar, navigate to **APIs & Services** > **Library**.
2. Search for and enable the following APIs:
   - **Google Sheets API**
   - **Google Drive API**

## Step 4: Create a Service Account
1. In the left sidebar, go to **IAM & Admin** > **Service Accounts**.
2. Click on **Create Service Account**.
3. Provide a name for the service account (e.g., `swiftcode-app-service-account`).
4. Under **Role**, choose **Project** > **Editor**.
5. Click **Create** and then **Done**.

## Step 5: Generate a JSON Key
1. In the Service Accounts list, locate the service account you just created.
2. Click on the three dots (menu) under **Actions** > **Create Key**.
3. Choose **JSON** as the key type and click **Create**.
4. A **`credentials.json`** file will be downloaded to your computer.

## Step 6: Place the Credentials File in the Project Directory
1. Move the downloaded **`credentials.json`** file into the root directory of your project.
   - The file should be in the same directory as your `docker-compose.yml` and `README.md`.

Your Google API credentials are now ready to be used by the SwiftCodeApp.
