# easy ui with tailwindcss

## dev

install tailwindcss cli

```bash
wget https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
chmod +x tailwindcss-linux-x64
mv tailwindcss-linux-x64 tailwindcsscli
```

Run the CLI tool to scan your source files for classes and build your CSS.
```bash
tailwindcsscli -i ./src/input.css -o ./src/output.css --watch
```