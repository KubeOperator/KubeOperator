import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {CredentialService} from '../credential.service';
import {Credential} from './credential';
import {HostInfoComponent} from '../../host/host-info/host-info.component';
import {TipLevels} from '../../tip/tipLevels';
import {TipService} from '../../tip/tip.service';

@Component({
  selector: 'app-credential-list',
  templateUrl: './credential-list.component.html',
  styleUrls: ['./credential-list.component.css']
})
export class CredentialListComponent implements OnInit {

  items: Credential[] = [];
  selected: Credential[] = [];
  loading = true;
  @Output() add = new EventEmitter();
  @ViewChild(HostInfoComponent, { static: false })
  child: HostInfoComponent;
  showDelete = false;
  resourceTypeName: '凭据';

  constructor(private credentialService: CredentialService, private tipService: TipService) {
  }

  ngOnInit() {
    this.listItems();
  }

  listItems() {
    this.loading = true;
    this.credentialService.listCredential().subscribe(data => {
      this.items = data;
      this.loading = false;
    });
  }

  delete() {
    const promises: Promise<{}>[] = [];
    this.selected.forEach(item => {
        promises.push(this.credentialService.deleteCredential(item.name).toPromise());
      }
    );
    Promise.all(promises).then(data => {
      this.tipService.showTip('删除成功', TipLevels.SUCCESS);
    }, error => {
      this.tipService.showTip('删除失败' + error.toString(), TipLevels.ERROR);
    }).finally(
      () => {
        this.showDelete = false;
        this.selected = [];
        this.listItems();
      }
    );
  }

  refresh() {
    this.listItems();
  }

  addItem() {
    this.add.emit();
  }

}
