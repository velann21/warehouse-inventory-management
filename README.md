# warehouse-inventory-management
This project is to manage the warehouse inventory.


Folder structures:
---
Deployments;
--
This will contains all the helm charts and kubernetes manifest files

Docs:
--
This is folder we will have all our documents related to this microservies

pkg:
--
The folder where we have all the project code base
    /controllers:
    ---
    Where the entry point for each API exist and will act as an controller for API functionality
    
    /helpers:
    ---
    Contains all helper functions
    
    /migration_scripts:
    ----
    Contains the migration scripts
    
    /models:
    ---
    This model has the all the model 
    1. request
    2.response
    3. database
    4.internals
    
    /repository
    ----
    Contains the database quries
    
    /routes:
    ---
    Contains all the route endpoints
    /service
    ---
    contaonis the business logic for an API
    
main.go
---
This is entry point/bootstrap for whole project  

Makefile:
---
Contains all the build commands

APIs:
----
1. curl --location --request POST 'http://localhost:8080/api/v1/inventory/articles' \
   --header 'Content-Type: application/json' \
   --data-raw '{
       "articles": [
           {
               "name": "srewdriver",
               "stock": "19"
           },
           {
               "name": "bolts",
               "stock": "19"
           },
           {
               "name": "bigger parts",
               "stock": "19"
           }
       ]
   }'
   
   Description:
   This API is to create bulk articles using body
   

2. curl --location --request POST 'http://localhost:8080/api/v1/inventory/articles/fromFile' \
   --form 'file=@/Users/singaravelannandakumar/Downloads/assignment/inventory.json'
   Description:
      This API is to create bulk articles using file
   
3. curl --location --request POST 'http://localhost:8080/api/v1/inventory/products' \
   --header 'Content-Type: application/json' \
   --data-raw '{
     "products": [
       {
         "name": "Dining Table",
         "contain_articles": [
           {
             "art_id": "24",
             "amount_of": "4"
           },
           {
             "art_id": "25",
             "amount_of": "8"
           },
           {
             "art_id": "26",
             "amount_of": "1"
           }
           ]
       }
     ]
   }
   '
   Description:
      This API is to create bulk products using body
   
4. curl --location --request POST 'http://localhost:8080/api/v1/inventory/products/fromFile' \
   --form 'file=@/Users/singaravelannandakumar/Downloads/assignment/products.json'
   Description:
      This API is to create bulk products using file
   
5. curl --location --request GET 'http://localhost:8080/api/v1/inventory/products'
   Description:
      This API is to list products

6. curl --location --request GET 'http://localhost:8080/api/v1/inventory/articles'
   Description:
         This API is to list articles

7. curl --location --request GET 'http://localhost:8080/api/v1/inventory/products/18'
   Description:
         This API is to get the product details of specific product

8. curl --location --request POST 'http://localhost:8080/api/v1/inventory/purchaseProducts' \
   --header 'Content-Type: application/json' \
   --data-raw '{
       "name": "Dinning Table",
       "id": "18"
   }'
   Description:
            This API is to buy the product based on productID and name

   
9. curl --location --request POST 'http://localhost:8080/api/v1/inventory/sqlmigration' \
   --header 'Content-Type: application/json' \
   --data-raw '{
       "upcount": 1,
       "downcount":0
   }' 
    Description:
               This API is to make an sql migration (To create tables) 
               
               
How to run the project:
---
I coudln't add these due to time constrains

This step I will explain during the our virtual interview/ Will show you how to run

Features to improve:
---
1. Need to add proper error handling
2. Need to create proper response structre for AddArticle API and AddProducts API with failed and success 
3. Need to use proper Transaction for Mysql db write and update
4. Make the app properly helmized for kubernetes support

Note:
--
This code is not production grade code, I have written this whole code in just 6 hours, So i couldn't focus on lots of good practises.
 I will make this code properly runnable in k8s and docker before our virtual interview
                         
     
