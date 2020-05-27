import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Credential} from '../credential';
import {CredentialService} from '../credential.service';

@Component({
    selector: 'app-credential-delete',
    templateUrl: './credential-delete.component.html',
    styleUrls: ['./credential-delete.component.css']
})
export class CredentialDeleteComponent implements OnInit {

    opened = false;
    items: Credential[] = [];
    @Output() deleted = new EventEmitter();

    constructor(private service: CredentialService) {
    }

    ngOnInit(): void {
    }


    open(items: Credential[]) {
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
        });
    }
}
