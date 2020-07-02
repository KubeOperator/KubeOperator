import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StorageProvisionerCreateNfsComponent } from './storage-provisioner-create-nfs.component';

describe('StorageProvisionerCreateNfsComponent', () => {
  let component: StorageProvisionerCreateNfsComponent;
  let fixture: ComponentFixture<StorageProvisionerCreateNfsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StorageProvisionerCreateNfsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StorageProvisionerCreateNfsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
