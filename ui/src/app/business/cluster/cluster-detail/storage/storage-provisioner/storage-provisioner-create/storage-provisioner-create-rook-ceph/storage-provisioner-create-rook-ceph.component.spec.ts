import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StorageProvisionerCreateRookCephComponent } from './storage-provisioner-create-rook-ceph.component';

describe('StorageProvisionerCreateRookCephComponent', () => {
  let component: StorageProvisionerCreateRookCephComponent;
  let fixture: ComponentFixture<StorageProvisionerCreateRookCephComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StorageProvisionerCreateRookCephComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StorageProvisionerCreateRookCephComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
