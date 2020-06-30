import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StorageClassCreateComponent } from './storage-class-create.component';

describe('StorageClassCreateComponent', () => {
  let component: StorageClassCreateComponent;
  let fixture: ComponentFixture<StorageClassCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StorageClassCreateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StorageClassCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
