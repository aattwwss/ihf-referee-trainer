document.getElementById('toggle-button').addEventListener('click', function () {
    const allQuestionCards = document.querySelectorAll('.question-card');
    const isShowingAnswers = this.className.includes('hide');

    allQuestionCards.forEach(card => {
        const correctAnswers = card.dataset.correct.split(',');
        const choices = card.querySelectorAll('.choice input[type="checkbox"]');

        choices.forEach(choice => {
            const parentLabel = choice.parentElement;
            if (!isShowingAnswers) {
                if (choice.checked && correctAnswers.includes(choice.value)) {
                    parentLabel.classList.add('correct');
                    parentLabel.classList.remove('wrong', 'missing');
                } else if (choice.checked && !correctAnswers.includes(choice.value)) {
                    parentLabel.classList.add('wrong');
                    parentLabel.classList.remove('correct', 'missing');
                } else if (!choice.checked && correctAnswers.includes(choice.value)) {
                    parentLabel.classList.add('missing');
                    parentLabel.classList.remove('correct', 'wrong');
                } else {
                    parentLabel.classList.remove('correct', 'wrong', 'missing');
                }
            } else {
                parentLabel.classList.remove('correct', 'wrong', 'missing');
            }
        });
    });

    if (isShowingAnswers) {
        this.innerHTML = '<i class="fas fa-eye"></i>';
        this.classList.remove('hide');
    } else {
        this.innerHTML = '<i class="fas fa-eye-slash"></i>';
        this.classList.add('hide');
    }
});

const CHOICE_CHECK_MAP_KEY = 'choiceCheckMap';
const READ_CHECK_MAP_KEY = 'readCheckMap';

// Load state from localStorage
function loadState() {
    const choiceCheckboxes = document.querySelectorAll('.choice input[type="checkbox"]');
    const choiceCheckMap = JSON.parse(localStorage.getItem(CHOICE_CHECK_MAP_KEY)) || {};
    choiceCheckboxes.forEach(checkbox => {
        const isChecked = choiceCheckMap[`${checkbox.name}-${checkbox.value}`];
        if (isChecked) {
            checkbox.checked = true;
        }
    });

    const readCheckboxes = document.querySelectorAll('.read-checkbox');
    const readCheckMap = JSON.parse(localStorage.getItem(READ_CHECK_MAP_KEY)) || {};
    readCheckboxes.forEach(checkbox => {
        const isChecked = readCheckMap[checkbox.closest('.question-card').id];
        if (isChecked) {
            checkbox.checked = true;
            checkbox.closest('.question-card').classList.add('read');
        }
    });
}

// Save state to localStorage
function saveState() {
    const choiceCheckboxes = document.querySelectorAll('.choice input[type="checkbox"]');
    const choiceCheckMap = {};
    choiceCheckboxes.forEach(checkbox => {
        choiceCheckMap[`${checkbox.name}-${checkbox.value}`] = checkbox.checked;
    });
    localStorage.setItem(CHOICE_CHECK_MAP_KEY, JSON.stringify(choiceCheckMap));

    const readCheckboxes = document.querySelectorAll('.read-checkbox');
    const readCheckMap = {};
    readCheckboxes.forEach(checkbox => {
        readCheckMap[checkbox.closest('.question-card').id] = checkbox.checked;
    });
    localStorage.setItem(READ_CHECK_MAP_KEY, JSON.stringify(readCheckMap));
}

document.addEventListener('DOMContentLoaded', loadState);

document.querySelectorAll('.choice input[type="checkbox"]').forEach(checkbox => {
    checkbox.addEventListener('change', saveState);
});

document.querySelectorAll('.read-checkbox').forEach(checkbox => {
    checkbox.addEventListener('change', function() {
        const card = this.closest('.question-card');
        if (this.checked) {
            card.classList.add('read');
        } else {
            card.classList.remove('read');
        }
        saveState();
    });
});

document.getElementById('menu-button').addEventListener('click', function () {
    const menuItems = document.getElementById('menu-items');
    menuItems.style.display = menuItems.style.display === 'block' ? 'none' : 'block';
    if (menuItems.style.display === 'block') {
        menuItems.focus();
    }
});

document.addEventListener('click', function(event) {
    const menuItems = document.getElementById('menu-items');
    const menuButton = document.getElementById('menu-button');
    if (menuItems.style.display === 'block' && !menuButton.contains(event.target) && !menuItems.contains(event.target)) {
        menuItems.style.display = 'none';
    }
});

function scrollToQuestion(questionId) {
    document.getElementById(questionId).scrollIntoView({ behavior: 'smooth' });
    document.getElementById('menu-items').style.display = 'none';
}
