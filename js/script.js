let listInputs = document.getElementsByClassName("list-inputs")[0];
const textWidths = {};
const shareButton = document.getElementById("share");
const editButton = document.getElementById("edit");


const wishlistItems = [];

const editList = () => {
  shareButton.classList.remove("hidden");
  editButton.classList.add("hidden");
  capDiv.classList.add("gift-cap-opened");

  createInputElement(listInputs.children.length);
};

const onFormSubmit = () => {
  capDiv.classList.remove("gift-cap-opened");
  shareButton.classList.add("hidden");
  editButton.classList.remove("hidden");

}

const createInputElement = (index) => {
  const inputRow = document.createElement("div");
  inputRow.setAttribute("class", "input-row");

  const inputContainer = document.createElement("div");
  inputContainer.setAttribute("class", "input-container");
  const input = document.createElement("input");
  input.setAttribute("placeholder", "Start typing");
  input.setAttribute("name", "item" + index);
  const indexSpan = document.createElement("span");
  indexSpan.innerText = "" + (index + 1) + ".";
  indexSpan.classList.add("invisible");
  inputContainer.appendChild(input);
  inputRow.appendChild(indexSpan);
  inputRow.appendChild(inputContainer);

  const span = document.createElement("span");
  span.style.visibility = "hidden";
  span.style.whiteSpace = "pre";
  span.style.position = "absolute";

  textWidths[index] = 0;

  input.oninput = (evt) => {
    // Set the span's text to the input's value
    span.textContent = input.value || input.placeholder;

    // Get the width of the span
    const textWidth = span.offsetWidth;
    textWidths[index] = textWidth;

    const maxWidth = Math.max(...Object.values(textWidths));
    const targetWidth = Math.min(Math.max(maxWidth + 140, 250), 400);
    listContainer.style.width = targetWidth + "px";
    capDiv.style.width = targetWidth + 70 + "px";

    const textValue = evt.target.value;
    if (textValue) {
      if (index === wishlistItems.length) {
        wishlistItems.push(textValue);
      } else {
        wishlistItems[index] = textValue;
      }
      indexSpan.classList.remove("invisible");
      if (listInputs.children.length === index + 1) {
        createInputElement(index + 1);
      }
    } else {
      indexSpan.classList.add("invisible");
      if (listInputs.children.length === index + 2) {
        const lastChild =
          listInputs.children[listInputs.children.length - 1];
        const lastChildValue = lastChild.value;

        if (!lastChildValue) {
          listInputs.removeChild(lastChild);
        }
      }
    }
  };

  input.addEventListener("keypress", function(event) {
    if (event.key === "Enter") {
      event.preventDefault();

      const nextInput = document.getElementsByName("item" + (index + 1))[0];
      if (nextInput) {
        nextInput.focus();
      }
    }
  });

  listInputs.appendChild(inputRow);
  document.body.appendChild(span);
};

// createInputElement(0);x§
