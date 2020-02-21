import {User} from '../user/user';

export class ItemMemberWrite {
  name: string;
  users: string[] = [];
}

export class ItemMember {
  name: string;
  users: User[] = [];
}
