// populates character selection list
function ui_populateChars(chars){
    var select = document.querySelector(".char-selection");
				
    chars.forEach(char => {
        var option = document.createElement("option");
        option.setAttribute("value", char);
        option.textContent = char;
        select.append(option)
    });
}

// deletes an element by a css class
function deleteCSSElementByClass(cssClass) {
    var elements = document.getElementsByClassName(cssClass);
    while(elements.length > 0){
        elements[0].parentNode.removeChild(elements[0]);
    }
}

// sets current value to maximum possible value
function updateItemToMaxQty() {
    i = event.target.id;
    console.log(i);
    var inputs = document.getElementsByClassName("count-input");
    var max_counts = document.getElementsByClassName("item-maxcount"); 
    inputs[i].value = Number.isInteger(parseInt(max_counts[i].textContent)) ?max_counts[i].textContent : 1 ;
}

function minmax(value, min, max) 
{
    if(parseInt(value) < min || isNaN(parseInt(value))) 
        return min; 
    else if(parseInt(value) > max) 
        return max; 
    else return value;
}