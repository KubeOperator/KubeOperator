import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {User} from '../user';
import {UserService} from '../user.service';

@Component({
  selector: 'app-user-list',
  templateUrl: './user-list.component.html',
  styleUrls: ['./user-list.component.css']
})
export class UserListComponent implements OnInit {

  loading = true;
  users: User[] = [];
  selectedRow: User[] = [];
  @Output() addUser = new EventEmitter();

  constructor(private userService: UserService) {
  }

  ngOnInit() {
    this.listUser();
  }

  listUser() {
    this.userService.listUsers().subscribe(data => {
      this.users = data;
      this.loading = false;
    });
  }

  addNewUser() {
    this.addUser.emit();
  }

  refresh() {
    this.listUser();
  }
}
