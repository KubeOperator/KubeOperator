import index from '@angular/cli/lib/cli';
import {Item} from '../item/item';

export class Profile {
  id: string;
  current_item: string;
  items: Item[] = [];
  user: SessionUser;
  item_role_mappings: ItemRoleMapping[] = [];
}

export class SessionUser {
  id: number;
  username: string;
  email: string;
  is_superuser: boolean;
  is_active: boolean;
  token: string;
  password: string;
  current_item: string;
}

export class ItemRoleMapping {
  role: string;
  item_name: string;
}
