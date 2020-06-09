export class SessionUser {
    userId: string;
    name: string;
    token: string;
    isActive: boolean;
    email: string;
    language: string;
}

export class Profile {
    user: SessionUser;
    token: string;
}
