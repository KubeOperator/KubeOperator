import index from '@angular/cli/lib/cli';

export class SessionUser {
  id: number;
  username: string;
  email: string;
  is_superuser: boolean;
  is_active: boolean;
  token: string;
  password: string;
}
