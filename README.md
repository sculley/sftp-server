# SFTP Server

<br />
<h3 align="center">SFTP Server</h3>

  <p align="center">
    sftp-server is an in-memory SFTP Server implementation written in Go that can be used in Go unit-tests to test code that interacts with an SFTP server.
    <br />
    <a href="https://github.com/sculley/sftp-server"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/sculley/sftp-server/issues">Report Bug/Issue</a>
    ·
    <a href="https://github.com/sculley/sftp-server/pulls">Request Feature</a>
  </p>
</div>

<summary>Table of Contents</summary>
<ol>
  <li>
    <a href="#about-the-project">About The Project</a>
    <ul>
      <li><a href="#built-with">Built With</a></li>
    </ul>
  </li>
  <li>
    <a href="#getting-started">Getting Started</a>
    <ul>
      <li><a href="#installation">Installation</a></li>
    </ul>
  </li>
  <li><a href="#usage">Usage</a></li>
  <li><a href="#roadmap">Roadmap</a></li>
  <li><a href="#contributing">Contributing</a></li>
  <li><a href="#license">License</a></li>
  <li><a href="#contact">Contact</a></li>
</ol>

## About The Project

sftp-server is an in-memory SFTP Server implementation written in Go that can be used in Go unit-tests to test code that interacts with an SFTP server. It is useful for ensuring that your code works correctly with an SFTP server without having to run an SFTP along with your tests or to have to mock the SFTP server which can be a complex task to do.

### Built With

* [![Go][Go-Badge]][Go-url]

## Getting Started

### Usage

```go
server := sftpserver.New("localhost:2022", "sculley", "password")
if err := server.Start(); err != nil {
    log.Fatal("Failed to start test SFTP server:", err)
}
defer server.Stop()

conn, err := ssh.Dial("tcp", "localhost:2022", &ssh.ClientConfig{
    User: "sculley",
    Auth: []ssh.AuthMethod{
        ssh.Password("password"),
    },
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
})
if err != nil {
    log.Fatal("Failed to dial SFTP server:", err)
}

client, err := sftp.NewClient(conn)
if err != nil {
    log.Fatal("Failed to create SFTP client:", err)
}

files, err := client.ReadDir("/tmp")
if err != nil {
    log.Fatal("Failed to list files: ", err)
}

for _, file := range files {
    log.Println(file.Name())
}
```

## Issues

See the [open issues](https://github.com/sculley/sftp-server/issues) for a full list of proposed features (and known issues).

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Clone or Fork the project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request to merge into develop

## License

Distributed under the Apache License. See `LICENSE.txt` for more information.

## Contact

Project Link: [https://github.com/sculley/sftp-server](https://github.com/sculley/sftp-server)

[Go-Badge]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev
