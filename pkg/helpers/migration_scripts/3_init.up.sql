create table product_inventory(product_id int NOT NULL, inventory_id int NOT NULL, quantity_each int, total_required_quantity int, PRIMARY KEY (product_id, inventory_id), FOREIGN KEY (product_id) REFERENCES products(id), FOREIGN KEY (inventory_id) REFERENCES inventory (art_id));