
<p align="center">
  <img src="https://github.com/user-attachments/assets/de5cc364-ddf5-44d2-a7bb-4e5b3bd85541" width="250" />
  <H1 align="center">Spread</H1>
</p>

**Spread** is an OTA update server for spreading new releases to React Native apps

## Usage 
The entire binary is a CLI tool capable of both running the server as well as supporting commands to release new versions. 

### Install Spread 
Install Spread to your system to interact with your remote Spread server:  
```sh
curl -fsSL https://cdn-swish.justswish.in/spread-install.sh | sh
```

### Make new release 

```sh
spread release \
  --remote REMOTE_SPREAD_HOST \
  --auth-key AUTH_KEY \
  --app-name APP_NAME \ # siwsh-ios
  --environment ENV \ # development, production
  --target-version VERSION \ # 1.0.1
  --os-name OS \ # ios, android
  --project-dir REACT_NATIVE_PROJECT_DIRECTORY \ 
  --is-typescript true \
  --description DESCRIPTION "fix: add to cart bug"
```
### Build new binary

```sh
make build
```
