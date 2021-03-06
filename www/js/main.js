// Connect UI actions to Go functions
const btnItemsSave = document.querySelector('.items-save');
const btnItemsGet = document.querySelector('.items-get');

// Action to run when "Save Data" button is clicked
btnItemsSave.addEventListener('click', () => {
    // get all item inputs
    var list_of_item_inputs = document.getElementsByClassName("count-input")
    var list_of_item_lvl_inputs = document.getElementsByClassName("lvl-input")

    // create a list with all input values sorted
    var int_qty_array = []
    for (let i = 0; i < list_of_item_inputs.length; i++) {
        const el = list_of_item_inputs[i];
        int_qty_array.push(parseInt(el.value));
    }

    // create a list with all input values sorted
    var int_lvl_array = []
    for (let i = 0; i < list_of_item_lvl_inputs.length; i++) {
        const el = list_of_item_lvl_inputs[i];
        int_lvl_array.push(parseInt(el.value));
    }

    // update items
    goUpdateItems(JSON.stringify(int_qty_array), JSON.stringify(int_lvl_array)).then(()=>{
        // update power
        var power_selection = document.getElementById( "power-selection" );
        goUpdatePower(power_selection.options[power_selection.selectedIndex].value).then( () => {
            // then save data
            goSaveData().then(()=>{
                alert('Saved Data!');
                loadItems();
            });
        });  
    });
});

// Action to run when "Load Items" button is clicked
btnItemsGet.addEventListener('click',  () => {
    loadItems();
});

function loadItems() {
    // get Character Items
    goGetItems().then(jsonItems => {
        // Get Character Powers
        goGetPowers().then(powers => {
            let jsonPowers = JSON.parse(powers);
    
            // clear powers
            let select = document.querySelector(".power-selection");
            deleteCSSElementByClass("power-option");
        
            // toggle visibility of equiped power selection
            let el = document.getElementById( "power-selection-container" );
            el.classList.add("power-selection-hidden");
            el.classList.remove("power-selection-hidden");
            
            jsonPowers.forEach(power => {
                let option = document.createElement("option");
                option.setAttribute("value", power);
                option.classList.add("power-option");
                option.textContent = power;
                select.append(option)
            });
        
            listRender(jsonItems);
        });
    });
}

// Render and Populate Chars list 
const render = () => {
    goGetChars().then(chars => {
        ui_populateChars(chars);
    });
}

render();


