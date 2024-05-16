db = db.getSiblingDB('zenith')
db.createUser(
    {
        user: "zenith",
        pwd: "zenith",
        roles: [
            {
                role: "readWrite",
                db: "zenith"
            }
        ]
    }
);