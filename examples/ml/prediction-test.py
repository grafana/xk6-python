"""
Example k6 test script to demonstrate the use of testing an LLM API

In this case we are testing an OLLAMA server. See the install instructions on https://github.com/ollama/ollama

Run:
- locally with dashboard: K6_WEB_DASHBOARD=true ./k6 run examples/ml/prediction-test.star
- against cloud: ./k6 run --out cloud examples/ml/prediction-test.star 
"""

load("requests", "post")
load("check", "check")
load("metrics", "gauge", "counter")

evalDurationGauge = gauge("eval_duration_seconds")
correctAnswers = counter("correct_answers_count")

options = {
 "stages": [
    { "duration": '30s', "target": 1 },
    { "duration": '30s', "target": 20 },
    { "duration": '30s', "target": 30 },
  ],
  "thresholds": {
    "checks": ["rate>=0.99"],
  },
}

def default(_):
    prompt = "Answer one word, yes or no: do people in %s like hotdogs"
    countries = {
        "the Netherlands": "Yes",
        "Spain": "No",
        "Belgium": "Yes",
        "The US": "Yes",
        "Hungary": "No",
        "Ireland": "Yes"
    }
    model = "llama2"

    for country, reply in countries.items():
        resp = post(
            'http://10.80.20.53:11434/api/generate',
            json = {
                "model": model,
                "prompt": prompt % country,
                "stream": False,
        })

        check(resp, {
            "is status 200": lambda r: r.status_code == 200,
        })

        if resp.ok:
            data = resp.json()
            # The evaluation duration as seen by the model
            evalDurationGauge.add(int(data['eval_duration'] / 1000), {"country": country, "answer": reply})

            # Did the model reply correctly
            correctAnswers.add(data['response'].find(reply) != -1, {"country": country, "answer": reply})
