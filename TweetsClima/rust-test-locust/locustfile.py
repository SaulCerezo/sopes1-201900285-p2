from locust import HttpUser, task, between
import json
import random

class RustApiUser(HttpUser):
    wait_time = between(1, 2)  # Espera entre peticiones

    @task
    def enviar_tweet(self):
        payload = [
            {
                "description": random.choice(["Soleado", "Nublado", "Lluvioso"]),
                "country": "GT",
                "weather": random.choice(["Soleado", "Nublado", "Lluvia"])
            }
        ]
        self.client.post("/input", data=json.dumps(payload), headers={"Content-Type": "application/json"})
