import json
import socket
import time
import logging

CONFIG_FILE = "config.json"
LOG_FILE = "log.txt"
BUF_SIZE = 1024

def read_config():
    with open(CONFIG_FILE, "r") as f:
        config = json.load(f)
    server_address = config.get("server_address")
    if not server_address:
        raise ValueError("server_address is empty")
    return server_address

def main():
    logging.basicConfig(filename=LOG_FILE, level=logging.INFO, format="%(asctime)s: %(message)s")
    try:
        server_address = read_config()
    except Exception as e:
        logging.error(f"Failed to read config: {e}")
        return

    try:
        host, port = server_address.split(":")
        port = int(port)
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.connect((host, port))
        logging.info(f"Connect to server {server_address}")
        time.sleep(2)
        
        msg = "Вовк Илья Богданович"
        sock.sendall(msg.encode())
        logging.info(f"Send message: {msg}")

        buf = sock.recv(BUF_SIZE)
        logging.info(f"Recieve message: {buf.decode(errors='replace')}")
        sock.recv(BUF_SIZE)

    except Exception as e:
        logging.error(f"Socket error: {e}")
    finally:
        sock.close()

if __name__ == "__main__":
    main()
