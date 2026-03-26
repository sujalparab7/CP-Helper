# AI-Driven Competitive Programming Advisor

A powerful, full-stack pedagogical tool designed to deeply analyze your Codeforces submissions, detect distinct algorithmic weaknesses, profile your behavioral cadence under pressure, and utilize Google's Gemini AI to synthesize dynamic Socratic hints alongside a highly personalized 5-day training matrix.

## Features

- **Automated Failure Tracking:** Directly hooks into the public Codeforces API to aggressively filter your most recent contest submissions and isolate your specific programmatic flaws.
- **Dynamic Behavioral Profiling:** Mathematically tracks the time-delta between your consecutive problem submissions to diagnose "panic-submitting" and "rushing" versus systematic, methodical dry-running.
- **5-Day Training Matrices:** Analyzes your historical algorithmic weak points (e.g., `greedy`, `dp`, `graphs`) and leverages the AI backend to automatically generate a tailored, actionable 5-day training schedule.
- **AI Socratic Tutor:** Instead of giving you the answers to puzzles, the platform dynamically scrapes the actual raw source code of your failed Codeforces submissions, passing it securely to the Gemini LLM. It returns 4 structured levels of pedagogical scaffolding (Conceptual, Structural, Diagnostic, and Edge-Case definitions).
- **Graceful Cascade Fallbacks:** Built-in multi-model exponential backoff scaling that guarantees rapid structural analysis even against Google API 503 limits, gracefully defaulting to safe metrics if system networks are offline.
- **Glassmorphic UI:** A breathtaking, minimalist frontend powered by raw HTML5/Vanilla CSS utilizing native expanding UI semantics and dynamic caching.

## Prerequisites

- **Golang** (v1.18 or higher)
- **Google Gemini API Key** (Accessible via Google AI Studio)

## Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/cp-advisor.git
   cd cp-advisor
   ```

2. **Configure your API Key:**
   Create a `.env` file in the root directory (same level as `main.go`) and insert your real Gemini API key:
   ```env
   GEMINI_API_KEY=your_gemini_api_key_here
   ```

3. **Run the backend server:**
   ```bash
   go run main.go
   ```

4. **Access the application:**
   Open your browser and navigate to `http://localhost:8080`.

## Tech Stack
- **Backend:** Golang (`net/http` standard library)
- **Frontend:** Vanilla HTML, CSS, JavaScript (No heavy frameworks)
- **External Integrations:** Codeforces Official API, Google Gemini AI (via REST)

## License
MIT License
