export class Credential {
  id: string;
  name: string;
  username: string;
  password: string;
  private_key: string;
  date_created: string;
  type = 'password';
}
