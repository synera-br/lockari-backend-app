# How to use the new tenant creation feature

This document explains how to use the new tenant creation feature, which is integrated with OpenFGA.

## Overview

When a new user signs up for the application, a new tenant is automatically created for them in OpenFGA. The user is also assigned as the administrator of this new tenant.

## How it works

The signup process now includes the following steps:

1.  A new user signs up.
2.  The `SignupEventService` is called.
3.  The service generates a new tenant ID.
4.  The service calls the `SetupTenant` method of the `LockariAuthorizationService`.
5.  The `LockariAuthorizationService` creates the new tenant in OpenFGA and assigns the user as the owner and administrator.
6.  The tenant ID is saved in Firebase Authentication for the user.
7.  The signup event is saved to the database.

## API Usage

There are no changes to the public API for the signup process. The tenant creation is handled automatically on the backend.

## Configuration

To use this feature, you need to configure the OpenFGA connection in your configuration file. Add the following section to your configuration file:

```json
{
  "openfga": {
    "api_url": "http://localhost:8080",
    "store_id": "YOUR_STORE_ID",
    "authorization_model_id": "YOUR_AUTHORIZATION_MODEL_ID"
  }
}
```

Replace the values with your OpenFGA API URL, store ID, and authorization model ID.
