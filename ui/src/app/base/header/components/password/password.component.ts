import {Component, OnInit, ViewChild} from '@angular/core';
import {SessionService} from '../../../../shared/session.service';
import {NgForm} from '@angular/forms';

@Component({
  selector: 'app-password',
  templateUrl: './password.component.html',
  styleUrls: ['./password.component.css']
})
export class PasswordComponent implements OnInit {

  oldPassword: string;
  newPassword: string;
  confirmPassword: string;

  constructor(private sessionService: SessionService) {
  }

  opened = false;
  submitGoing = false;
  alertOpen = false;
  msg = '';
  level = '';

  @ViewChild('passform', {static: true}) passform: NgForm;

  ngOnInit() {
  }

  checkOldPassword() {
    const user = this.sessionService.getCacheUser();

    this.sessionService.changePassword(user.id, this.oldPassword, this.newPassword).subscribe(data => {
      this.showMsg('success', '修改成功');
      this.onCancel();
    }, error => {
      this.showMsg('danger', '修改失败,原密码错误!');
      this.clear();
    });
  }

  onSubmit() {
    if (this.submitGoing) {
      return;
    }
    if (this.newPassword !== this.confirmPassword) {
      this.showMsg('danger', '两次密码输入不一致!');
      this.clear();
    } else {
      this.checkOldPassword();
    }

  }

  clear() {
    this.oldPassword = null;
    this.newPassword = null;
    this.confirmPassword = null;
  }

  onCancel() {
    this.opened = false;
    this.clear();
    this.msg = '';
    this.alertOpen = false;
    this.level = '';
    this.passform.resetForm();
  }

  showMsg(level, msg) {
    this.level = level;
    this.msg = msg;
    this.alertOpen = true;
    setInterval(() => {
      this.alertOpen = false;
    }, 2000);
  }


}
