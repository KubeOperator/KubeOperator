export class Profile {
    user: SessionUser;
    token: string;
}

export class SessionUser {
    userId: string;
    name: string;
    token: string;
    isActive: string;
    email: string;
    language: string;
}