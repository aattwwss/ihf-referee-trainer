const multiSelectWithoutCtrl = (elemSelector) => {
    let options = [].slice.call(document.querySelectorAll(`${elemSelector} option`));
    options.forEach(function (element) {
        element.addEventListener("mousedown", function (e) {
            e.preventDefault();
            element.parentElement.focus();
            this.selected = !this.selected;
            return false;
        }, false);
    });
}

document.addEventListener('DOMContentLoaded', function () {
    multiSelectWithoutCtrl('#rules');
});

// document.getElementById('num-questions').addEventListener('input', function (e) {
//     e.target.value = e.target.value.replace(/[^0-9.]/g, '');
//     const numQuestions = parseInt(e.target.value);
//     const maxNumQuestions = parseInt(e.target.max);
//     if (numQuestions > maxNumQuestions) {
//         e.target.value = maxNumQuestions;
//     }
// })

const textInputs = Array.from(document.getElementsByClassName('number-input'));
textInputs.forEach(el => el.addEventListener('input', function (e) {
        e.target.value = e.target.value.replace(/[^0-9.]/g, '');
        const numQuestions = parseInt(e.target.value);
        const maxNumQuestions = parseInt(e.target.max);
        if (!isNaN(maxNumQuestions) && numQuestions > maxNumQuestions) {
            e.target.value = maxNumQuestions;
        }
    })
)