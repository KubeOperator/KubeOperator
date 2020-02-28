import {Component, OnInit, ViewChild} from '@angular/core';
import {SessionService} from '../../../../shared/session.service';
import {NgForm} from '@angular/forms';

@Component({
  selector: 'app-password',
  templateUrl: './password.component.html',
  styleUrls: ['./password.component.css']
})
export class PasswordComponent implements OnInit {

  original: string;
  password: string;
  confirmPassword: string;

  constructor(private sessionService: SessionService) {
  }

  opened = false;
  submitGoing = false;

  @ViewChild('passform', {static: true}) passform: NgForm;

  ngOnInit() {
  }

  onSubmit() {
    if (this.submitGoing) {
      return;
    }
    this.submitGoing = true;
    if (this.password !== this.confirmPassword) {
      this.submitGoing = false;
      return;
    } else {
      this.sessionService.changePassword(this.original, this.password).subscribe(data => {
      }, error => {
        console.log(error);
      });
    }
  }

  checkPassword() {
    return this.password === this.confirmPassword;
  }

  clear() {
    this.passform.resetForm();
    this.original = null;
    this.password = null;
    this.confirmPassword = null;
  }

  onCancel() {
    this.opened = false;
    this.clear();
  }
}
