document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('analyzeForm');
    const handleInput = document.getElementById('handleInput');
    const loading = document.getElementById('loading');
    const errorDiv = document.getElementById('error');
    const resultsSection = document.getElementById('resultsSection');

    const profileContent = document.getElementById('profileContent');
    const weaknessesList = document.getElementById('weaknessesList');
    const trainingMatrix = document.querySelector('#trainingMatrix tbody');
    const aiFeedbackContext = document.getElementById('aiFeedbackContext');
    const aiFeedbackContent = document.getElementById('aiFeedbackContent');

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const handle = handleInput.value.trim();
        if (!handle) return;

        // Reset UI state
        resultsSection.classList.add('hidden');
        errorDiv.classList.add('hidden');
        loading.classList.remove('hidden');

        try {
            const baseUrl = window.location.protocol === 'file:' ? 'http://localhost:8080' : '';
            const response = await fetch(`${baseUrl}/api/analyze?handle=${encodeURIComponent(handle)}`);
            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'Failed to fetch analysis from backend');
            }

            renderDashboard(data);
        } catch (err) {
            errorDiv.textContent = err.message;
            errorDiv.classList.remove('hidden');
        } finally {
            loading.classList.add('hidden');
        }
    });

    function renderDashboard(data) {
        profileContent.innerHTML = `
            <p style="margin-bottom: 0.5rem"><strong>Overall Status:</strong> ${data.profile.status}</p>
            <p style="margin-bottom: 0.5rem"><strong>Submission Cadence:</strong> <span style="color:var(--accent-color)">${data.profile.cadence}</span></p>
            <p><strong>Diagnosis:</strong> <span style="color:var(--text-secondary)">${data.profile.notes}</span></p>
        `;

        weaknessesList.innerHTML = '';
        if (data.weaknesses && data.weaknesses.length > 0) {
            data.weaknesses.forEach(w => {
                const li = document.createElement('li');
                li.className = 'tag';
                li.textContent = `${w.tag} (${w.count} fails)`;
                weaknessesList.appendChild(li);
            });
        } else {
            weaknessesList.innerHTML = '<li class="tag">No major weaknesses identified yet.</li>';
        }

        trainingMatrix.innerHTML = '';
        if (data.matrix && data.matrix.length > 0) {
            data.matrix.forEach(day => {
                const tr = document.createElement('tr');
                tr.innerHTML = `
                    <td><strong>${day.day}</strong></td>
                    <td class="focus-col">${day.focus || ''}</td>
                    <td>${day.objective}</td>
                    <td>${day.action}</td>
                `;
                trainingMatrix.appendChild(tr);
            });
        }

        aiFeedbackContent.innerHTML = '';
        if (data.aiFeedback) {
            aiFeedbackContext.innerHTML = `<strong>Context:</strong> Failed submission on Problem <em>${data.aiFeedback.problemName}</em>`;
            
            const levels = [
                { level: 1, title: 'Concept Abstraction', content: data.aiFeedback.level1 },
                { level: 2, title: 'Structural Guidance', content: data.aiFeedback.level2 },
                { level: 3, title: 'Diagnostic Feedback', content: data.aiFeedback.level3 },
                { level: 4, title: 'Test-Case Instantiation', content: data.aiFeedback.level4 }
            ];

            levels.forEach(lvl => {
                const details = document.createElement('details');
                details.className = `feedback-level level-${lvl.level}`;
                details.innerHTML = `
                    <summary>Level ${lvl.level}: ${lvl.title}</summary>
                    <div class="details-content">
                        <p>${lvl.content}</p>
                    </div>
                `;
                aiFeedbackContent.appendChild(details);
            });
        }

        resultsSection.classList.remove('hidden');
    }
});
