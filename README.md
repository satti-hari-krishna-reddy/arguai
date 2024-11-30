# Argu AI - Debate Platform

**Argu AI** is an AI-powered debate platform built for the **ModusHack Hackathon**. It allows users to choose debate topics, assign AI personas, and watch the models engage in real-time debates.

This project leverages **Modus** for orchestrating AI models and **Go** as the programming language, making it an interactive and real-time AI debate experience.

---

## 🚀 Project Overview

Argu AI is designed to facilitate engaging and interactive AI debates. Users can select two AI models from multiple options (e.g., Gemini, GPT-4, Meta LLaMA) and assign them different personas. The debate occurs in a back-and-forth format, where the models argue different viewpoints on a given topic. After 8 rounds of debate, a third AI model acts as the **judge** to decide the winner based on the debate performance.


---

## 💻 Technical Details

### Modus for Model Orchestration

**Modus** is used to seamlessly orchestrate the communication between multiple AI models. The output from one model is fed as input into the next, creating a fluid and continuous debate process.

### Real-Time Data Storage

Debate messages arestreamed and stored in PostgreSQL (via Azure Cosmos DB) for real-time updates. This ensures that frontend can track the ongoing debate without significant delays.

---

## 📸 Debate Flow

Here’s how the debate works internally:

1. **User Inputs Topic**: The user selects a topic for the debate.
2. **Model Selection**: Two AI models are selected, each with a unique persona.
3. **Debate Flow**: The models debate in real-time, alternating their arguments.
4. **AI Judge Evaluation**: After 8 rounds, a third AI model judges the debate.
5. **Winner Announcement**: The winning model is decided based on the judge's evaluation.

![Debate Flow](path-to-your-image-1.png)

---

### ⚙️ Setup AI Models
Note: For the hackathon, the available models are limited due to API rate limits. Please ensure to use the Gemini model as it has higher rate limits..

## ⚡ Getting Started

To run this project locally, follow these steps:

### 1. Prerequisites

Before running the project, ensure you have the following installed:

- **Go** (for backend development)
-  Install the **Modus CLI** globally with npm:

  ```
  npm install -g @hypermode/modus-cli
```

### Cloning the Repository

```
git clone https://github.com/yourusername/argu-ai.git

cd arguai
```
Next, replace the environment variables in `.env.dev.local` with your credentials:
```
MODUS_OPENAI_API_KEY=
MODUS_GOOGLE_API_KEY=
MODUS_CONVDB_PG_USER=
MODUS_CONVDB_PG_PASSWORD=
```


### Run the App Locally
To run the app with fast refresh, use the following command:
```
 modus dev
```
This will start your app at http://localhost:8686/graphql.