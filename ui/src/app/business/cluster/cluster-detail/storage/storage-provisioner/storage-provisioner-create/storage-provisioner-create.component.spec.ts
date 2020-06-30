import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StorageProvisionerCreateComponent } from './storage-provisioner-create.component';

describe('StorageProvisionerCreateComponent', () => {
  let component: StorageProvisionerCreateComponent;
  let fixture: ComponentFixture<StorageProvisionerCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StorageProvisionerCreateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StorageProvisionerCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
