import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelComponent} from '../../../shared/class/BaseModelComponent';
import {Host} from '../host';
import {HostService} from '../host.service';

@Component({
    selector: 'app-host-delete',
    templateUrl: './host-delete.component.html',
    styleUrls: ['./host-delete.component.css']
})
export class HostDeleteComponent extends BaseModelComponent<Host> implements OnInit {

    opened = false;
    items: Host[] = [];
    @Output() deleted = new EventEmitter();

    constructor(private hostService: HostService) {
        super(hostService);
    }

    ngOnInit(): void {
    }

    open(items) {
        this.opened = true;
        this.items = items;
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
