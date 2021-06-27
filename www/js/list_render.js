async function listRender (jsonItems) {
    try {
        // delete existing options before populating new ones
        deleteCSSElementByClass("item-row");
        deleteCSSElementByClass("item-header");

        var jsonObj = JSON.parse(jsonItems);
        var new_items = [];
    
        jsonObj.forEach( (element, i) => {
            let this_element = {
                id: element.FileItem.Name,
                name: element.DbItem.Name,
                weight: element.DbItem.Weight,
                lvl: element.FileItem.Lvl,
                count: element.FileItem.OriginalCount,
                index: element.FileItem.index,
                max_count: element.DbItem.Stack,
                teleportable: element.DbItem.Teleportable,
            };
            new_items.push(this_element);
        });

        // get table element
        var table = document.querySelector("table");
        table.classList.add("item-container")

        // add header row and header cells
        var header_row = document.createElement("tr");
        header_row.classList.add("item-header")

        var header_thumb = document.createElement("th");
        header_thumb.classList.add("item-header-thumb")
        header_thumb.textContent = "Item";

        var header_name = document.createElement("th");
        header_name.classList.add("item-header-name")
        header_name.textContent = "Name";

        var header_lvl = document.createElement("th");
        header_lvl.classList.add("item-header-lvl")
        header_lvl.textContent = "Lvl";

        var header_weight = document.createElement("th");
        header_weight.classList.add("item-header-weight")
        header_weight.textContent = "Weight";

        var header_count = document.createElement("th");
        header_count.classList.add("item-header-count")
        header_count.textContent = "Qty";

        var header_maxcount = document.createElement("th");
        header_maxcount.classList.add("item-header-maxcount")
        header_maxcount.textContent = "MaxQty";

        var header_actions = document.createElement("th");
        header_actions.classList.add("item-header-actions")
        header_actions.textContent = "Actions";

        header_row.append(header_thumb);
        header_row.append(header_name);
        header_row.append(header_lvl);
        header_row.append(header_weight);
        header_row.append(header_count);
        header_row.append(header_maxcount);
        header_row.append(header_actions);
        table.append(header_row);

        // add a row for each item in inventory
        for (var i = 0; i < new_items.length; i++) {
            var item = new_items[i];

            var item_row = document.createElement("tr");
            item_row.classList.add("item-row")

            var item_thumb = document.createElement("td");
            item_thumb.classList.add("item-thumb")
            var item_thumb_img = document.createElement("img");
            item_thumb_img.classList.add("item-thumb-img");
            item_thumb_img.setAttribute("src", "img/items/" + item.id + ".png");
            item_thumb.append(item_thumb_img);
            item_row.append(item_thumb);

            var item_name = document.createElement("td");
            item_name.classList.add("item-name")
            item_name.textContent = item.name;
            item_row.append(item_name);

            // cell for lvl
            var item_lvl_cell = document.createElement("td");
            item_lvl_cell.classList.add("item-lvl");
            
            // input for lvl
            var item_lvl = document.createElement("input");
            item_lvl.classList.add("lvl-input")
            item_lvl.classList.add("item-input")
            item_lvl.setAttribute("type", "number");
            item_lvl.setAttribute("min", "1");
            item_lvl.setAttribute("max", "255");
            item_lvl.setAttribute("maxlength", "3");
            item_lvl.setAttribute("onkeyup", "this.value = minmax(this.value, 1, 255)");

            item_lvl.value = item.lvl;
            item_lvl_cell.append(item_lvl);
            item_row.append(item_lvl_cell);

            // cell for weight
            var item_weight = document.createElement("td");
            item_weight.classList.add("item-weight")
            item_weight.textContent = item.weight;
            item_row.append(item_weight);

            // cell for qty
            var item_count_cell = document.createElement("td");
            item_count_cell.classList.add("item-count");
            // input for qty
            var item_count = document.createElement("input");
            item_count.classList.add("count-input");
            item_count.classList.add("item-input");
            item_count.setAttribute("type", "number");
            item_count.setAttribute("min", "1");
            let eval_max_qty = Number.isInteger(parseInt(item.max_count)) ? item.max_count : 1;
            item_count.setAttribute("max", eval_max_qty);
            item_count.setAttribute("onkeyup", "this.value = minmax(this.value, 1, "+eval_max_qty+")");
            item_count.setAttribute("maxlength", "3");
            item_count.value = item.count;
            item_count_cell.append(item_count);
            item_row.append(item_count_cell);

            // cell for max qty
            var item_maxcount = document.createElement("td");
            item_maxcount.classList.add("item-maxcount")
            item_maxcount.textContent = item.max_count;
            item_row.append(item_maxcount);

            var item_actions_cell = document.createElement("td");
            item_actions_cell.classList.add("item-maxcount-action")

            var max_action = document.createElement("div");
            max_action.classList.add("btn")
            max_action.classList.add("items-action-max")
            max_action.setAttribute("id", i);
            max_action.setAttribute("onclick", "updateItemToMaxQty()")
            max_action.textContent = "MAX";
            item_actions_cell.append(max_action);

            item_row.append(item_actions_cell);

            table.appendChild(item_row);
        }
        // items_string = JSON.stringify(new_items)
        // debug.innerText = `Debug: ${items_string}`;
    }
    catch (e) {
        console.log(e);
    }
};