# Go Microservices Project

Une application complète de microservices construite avec Go, présentant différents modèles et protocoles de communication.

## Aperçu du Projet

Ce projet démontre une architecture de microservices utilisant Go, avec des services communiquant à travers différents protocoles:

- **REST API** - Pour la gestion des utilisateurs
- **gRPC** - Pour le catalogue de produits
- **GraphQL** - Pour les avis sur les produits
- **Kafka** - Pour le traitement des commandes

Chaque service possède sa propre base de données et communique via une API Gateway qui fournit une interface unifiée pour les clients.

## Architecture

### Composants

1. **API Gateway**: Point d'entrée central qui route les requêtes vers les services appropriés
2. **User Service**: Service REST pour la gestion des utilisateurs
3. **Product Service**: Service gRPC pour la gestion du catalogue de produits
4. **Review Service**: Service GraphQL pour les avis sur les produits
5. **Order Service**: Service basé sur Kafka pour le traitement des commandes
6. **Bases de données**: Instances PostgreSQL pour chaque service
7. **Kafka & Zookeeper**: Infrastructure de streaming d'événements

## Détails des Services

| Service | Protocole | Port | Base de données | Description |
|---------|----------|------|----------|-------------|
| API Gateway | HTTP | 8080 | N/A | Route les requêtes vers les services appropriés et sert le frontend |
| User Service | REST | 8081 | postgres-user (userdb) | Gère les comptes utilisateurs et l'authentification |
| Product Service | gRPC | 8082 | postgres-product (productdb) | Gère les informations sur les produits et le catalogue |
| Review Service | GraphQL | 8083 | postgres-review (reviewdb) | Gère les avis et les évaluations de produits |
| Order Service | Kafka | N/A | postgres-order (orderdb) | Traite et gère les commandes |

## Installation et Configuration

### Prérequis

- Docker et Docker Compose
- Go (version 1.21 ou supérieure) - pour le développement local
- Git

### Exécution de l'Application

1. Cloner le dépôt:
   ```bash
   git clone https://github.com/yourusername/go-microservices-project.git
   cd go-microservices-project
   ```

2. Démarrer l'application avec Docker Compose:
   ```bash
   docker-compose up -d
   ```

3. Vérifier que les services sont en cours d'exécution:
   ```bash
   docker-compose ps
   ```

4. Accéder à l'application:
   - Frontend: http://localhost:8080
   - Vérification de l'état: http://localhost:8080/health

### Résolution des problèmes courants

Si vous rencontrez des erreurs liées aux ports déjà utilisés:

1. Vérifiez quels processus utilisent les ports requis:
   ```bash
   sudo lsof -i :8080
   sudo lsof -i :8081
   ```

2. Arrêtez ces processus ou modifiez les ports dans docker-compose.yml si nécessaire.

### Configuration de Développement

Pour le développement local sans conteneurs:

1. Installer les dépendances:
   ```bash
   go mod download
   ```

2. Exécuter l'infrastructure requise (bases de données, Kafka):
   ```bash
   docker-compose up -d postgres-user postgres-product postgres-review postgres-order zookeeper kafka
   ```

3. Exécuter les services individuellement:
   ```bash
   # Exécuter l'API Gateway
   cd api-gateway
   go run main.go

   # Exécuter le User Service
   cd services/rest-service
   go run cmd/main.go

   # Exécuter les autres services de façon similaire
   ```

## Documentation de l'API

### API du User Service (REST)

- **GET /api/users/** - Obtenir tous les utilisateurs
- **GET /api/users/{id}** - Obtenir un utilisateur par ID
- **POST /api/users/** - Créer un nouvel utilisateur
- **PUT /api/users/{id}** - Mettre à jour un utilisateur existant
- **DELETE /api/users/{id}** - Supprimer un utilisateur

### API du Product Service (gRPC)

- **GET /api/products/** - Obtenir tous les produits
- **GET /api/products/{id}** - Obtenir un produit par ID
- **POST /api/products/** - Créer un nouveau produit
- **PUT /api/products/{id}** - Mettre à jour un produit existant
- **DELETE /api/products/{id}** - Supprimer un produit

### API du Review Service (GraphQL)

- **GET /api/reviews/product/{productId}** - Obtenir les avis pour un produit
- **GET /api/reviews/{id}** - Obtenir un avis par ID
- **POST /api/reviews/** - Créer un nouvel avis

Également accessible via l'endpoint GraphQL:
- **POST /graphql** - Point d'entrée GraphQL avec le corps de requête:
  ```json
  {
    "query": "{ reviewsByProduct(productId: 1) { id rating comment user { username } } }"
  }
  ```

### API du Order Service (Kafka)

- **GET /api/orders/** - Obtenir toutes les commandes
- **GET /api/orders/{id}** - Obtenir une commande par ID
- **POST /api/orders/** - Créer une nouvelle commande
- **PUT /api/orders/{id}** - Mettre à jour une commande
- **GET /api/orders/status/{id}** - Obtenir le statut d'une commande

## Tests

Pour tester facilement toutes les API, vous pouvez utiliser la collection Postman fournie:

1. Importez le fichier `postman_collection.json` dans Postman
2. Exécutez les requêtes pour tester les différents endpoints

## Structure du Projet

```
go-microservices-project/
├── api-gateway/               # Service API Gateway
│   ├── config/                # Configuration
│   ├── handlers/              # Gestionnaires de requêtes pour chaque service
│   ├── routes/                # Définitions des routes
│   ├── Dockerfile             # Instructions Docker
│   └── main.go                # Point d'entrée
├── services/                  # Microservices
│   ├── rest-service/          # Service User REST
│   │   ├── cmd/               # Point d'entrée
│   │   ├── config/            # Configuration
│   │   ├── controllers/       # Gestionnaires de requêtes
│   │   ├── database/          # Connexion à la base de données
│   │   ├── models/            # Modèles de données
│   │   ├── repository/        # Couche d'accès aux données
│   │   ├── routes/            # Définitions des routes
│   │   ├── services/          # Logique métier
│   │   └── Dockerfile         # Instructions Docker
│   ├── grpc-service/          # Service Product gRPC
│   ├── graphql-service/       # Service Review GraphQL
│   └── kafka-service/         # Service Order Kafka
├── proto/                     # Définitions Protocol Buffer pour gRPC
├── frontend/                  # Fichiers HTML/JS du frontend
│   ├── index.html             # Page principale
│   ├── users.html             # Page des utilisateurs
│   ├── products.html          # Page des produits
│   ├── reviews.html           # Page des avis
│   ├── orders.html            # Page des commandes
│   └── static/                # Ressources statiques
├── docker-compose.yml         # Configuration Docker Compose
├── postman_collection.json    # Collection Postman pour tester les API
└── README.md                  # Documentation du projet
```

## Workflow de Développement

1. Créez une branche à partir de `develop` pour votre fonctionnalité
   ```bash
   git checkout develop
   git checkout -b feature/nom-fonctionnalite
   ```

2. Développez votre fonctionnalité et commitez régulièrement
   ```bash
   git add .
   git commit -m "Description des changements"
   ```

3. Poussez votre branche et créez une Pull Request vers `develop`
   ```bash
   git push origin feature/nom-fonctionnalite
   ```

4. Après révision et approbation, fusionnez dans `develop`

## Ressources d'Apprentissage

Pour approfondir vos connaissances sur les concepts utilisés dans ce projet:

- [Go Documentation](https://golang.org/doc/)
- [gRPC Documentation](https://grpc.io/docs/)
- [GraphQL Documentation](https://graphql.org/learn/)
- [Kafka Documentation](https://kafka.apache.org/documentation/)
- [Docker Documentation](https://docs.docker.com/)

## Contributeurs

Ce projet a été développé par Boussaid Mohamed Amine dans le cadre de l'examen du module SOA et Microservices, sous la direction du Dr. Salah Gontara, pour l'année universitaire 2024-2025.
