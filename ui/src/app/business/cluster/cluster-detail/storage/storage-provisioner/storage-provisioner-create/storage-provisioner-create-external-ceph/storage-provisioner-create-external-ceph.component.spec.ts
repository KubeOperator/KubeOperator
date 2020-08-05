import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StorageProvisionerCreateExternalCephComponent } from './storage-provisioner-create-external-ceph.component';

describe('StorageProvisionerCreateExternalCephComponent', () => {
  let component: StorageProvisionerCreateExternalCephComponent;
  let fixture: ComponentFixture<StorageProvisionerCreateExternalCephComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StorageProvisionerCreateExternalCephComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StorageProvisionerCreateExternalCephComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
