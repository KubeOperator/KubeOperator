import index from '@angular/cli/lib/cli';

export class Profile {
  id: string;
  current_item: string;
  user: SessionUser;
}

export class SessionUser {
  id: number;
  username: string;
  email: string;
  is_superuser: boolean;
  is_active: boolean;
  token: string;
  password: string;
  items: string[] = [];
  current_item: string;
}
