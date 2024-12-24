import numpy as np
from flask import Flask, request, jsonify
import pandas as pd
from sklearn.preprocessing import LabelEncoder
from sklearn.multioutput import MultiOutputClassifier
from sklearn.ensemble import RandomForestClassifier

app = Flask(__name__)

# Загрузка и подготовка модели
le_gender = LabelEncoder()
le_cuisine = LabelEncoder()
le_excluded = LabelEncoder()
le_goal = LabelEncoder()

model = MultiOutputClassifier(RandomForestClassifier(n_estimators=150, random_state=98))
df = pd.read_csv('training_sample.csv')
df['gender'] = le_gender.fit_transform(df['gender'])
df['preferred_cuisine'] = le_cuisine.fit_transform(df['preferred_cuisine'])
df['excluded_categories'] = le_excluded.fit_transform(df['excluded_categories'])
df['goal'] = le_goal.fit_transform(df['goal'])

X = df[['gender', 'height', 'age', 'weight', 'preferred_cuisine', 'excluded_categories', 'goal']]
y_breakfast = df[['breakfast_1', 'breakfast_2', 'breakfast_3']].fillna(0)
y_lunch = df[['lunch_1', 'lunch_2', 'lunch_3']].fillna(0)
y_dinner = df[['dinner_1', 'dinner_2', 'dinner_3']].fillna(0)
y = pd.concat([y_breakfast, y_lunch, y_dinner], axis=1)

model.fit(X, y)

import numpy as np


@app.route('/predict', methods=['POST'])
def predict():
    data = request.json
    user_input = pd.DataFrame([{
        'gender': le_gender.transform([data['gender']])[0],
        'height': data['height'],
        'age': data['age'],
        'weight': data['weight'],
        'preferred_cuisine': le_cuisine.transform([data['preferred_cuisine']])[0],
        'excluded_categories': le_excluded.transform([data['excluded_categories']])[0],
        'goal': le_goal.transform([data['goal']])[0]
    }])

    probabilities = model.predict_proba(user_input)

    def random_choice(probs, used_indices):
        available_indices = [i for i in range(len(probs)) if i not in used_indices]
        if not available_indices:
            return 0
        normalized_probs = probs[available_indices] / np.sum(probs[available_indices])
        return np.random.choice(available_indices, p=normalized_probs)

    def get_predictions(probabilities, start, end):
        used_indices = set()
        predictions = []
        for i in range(start, end):
            prediction = random_choice(probabilities[i][0], used_indices)
            predictions.append(prediction)
            used_indices.add(prediction)
        return predictions

    breakfast = get_predictions(probabilities, 0, 3)
    lunch = get_predictions(probabilities, 3, 6)
    dinner = get_predictions(probabilities, 6, 9)

    result = {
        "breakfast": [int(x) for x in breakfast],
        "lunch": [int(x) for x in lunch],
        "dinner": [int(x) for x in dinner]
    }
    print(result)
    return jsonify(result)


if __name__ == '__main__':
    app.run(debug=True)
