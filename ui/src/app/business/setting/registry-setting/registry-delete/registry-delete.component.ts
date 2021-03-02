import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Credential} from '../../credential/credential';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {Registry} from '../registry';
import {CredentialService} from '../../credential/credential.service';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {RegistryService} from '../registry.service';

@Component({
  selector: 'app-registry-delete',
  templateUrl: './registry-delete.component.html',
  styleUrls: ['./registry-delete.component.css']
})
export class RegistryDeleteComponent implements OnInit {
  opened = false;
  items: Registry[] = [];
  @Output() deleted = new EventEmitter();
  constructor(private service: RegistryService, private modalAlertService: ModalAlertService,
              private commonAlertService: CommonAlertService, private translateService: TranslateService) {
  }

  ngOnInit(): void {
  }


  open(items: Registry[]) {
    this.items = items;
    this.opened = true;
  }

  onCancel() {
    this.opened = false;
  }

  onSubmit() {
    this.service.batch('delete', this.items).subscribe(data => {
      this.deleted.emit();
      this.opened = false;
      this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
    }, error => {
      this.opened = false;
      this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
    });
  }

}
